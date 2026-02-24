package storage

import (
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_PutToken(t *testing.T) {
	homeDir := testSetupEHomeDir(t)
	tokenPath := filepath.Join(homeDir, parentDirName, storageDirName, tokenFileName)
	badTokenPath := filepath.Join(homeDir, "not_exist", tokenFileName)

	storage, err := NewFileSorage()
	require.Nil(t, err, "Create file storage")

	tests := []struct {
		name    string
		path    string
		wantErr string
	}{
		{name: "success", path: tokenPath},
		{
			name: "fail",
			path: badTokenPath,
			wantErr: "failed to write token: open " +
				badTokenPath + ": no such file or directory",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			storage.tokenPath = test.path

			err := storage.PutToken("token_string")
			var gotErr string
			if err != nil {
				gotErr = err.Error()
			}
			assert.Equal(t, test.wantErr, gotErr, "Put token error")
		})
	}
}

func Test_Token(t *testing.T) {
	homeDir := testSetupEHomeDir(t)
	tokenPath := filepath.Join(homeDir, parentDirName, storageDirName, tokenFileName)
	badTokenPath := filepath.Join(homeDir, "not_exist", tokenFileName)

	storage, err := NewFileSorage()
	require.Nil(t, err, "Create file storage")

	err = storage.PutToken("token_string")
	require.Nil(t, err, "Set token value")

	tests := []struct {
		name    string
		path    string
		want    string
		wantErr string
	}{
		{name: "success", path: tokenPath, want: "token_string"},
		{
			name: "fail",
			path: badTokenPath,
			wantErr: "failed to get token: open " +
				badTokenPath + ": no such file or directory",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			storage.tokenPath = test.path

			got, err := storage.Token()
			assert.Equal(t, test.want, got, "Get token")
			var gotErr string
			if err != nil {
				gotErr = err.Error()
			}
			assert.Equal(t, test.wantErr, gotErr, "Get token error")
		})
	}
}

func Test_PutKey(t *testing.T) {
	homeDir := testSetupEHomeDir(t)
	keyPath := filepath.Join(homeDir, parentDirName, storageDirName, keyFileName)
	badKeyPath := filepath.Join(homeDir, "not_exist", keyFileName)

	storage, err := NewFileSorage()
	require.Nil(t, err, "Create file storage")

	tests := []struct {
		name    string
		path    string
		wantErr string
	}{
		{name: "success", path: keyPath},
		{
			name: "fail",
			path: badKeyPath,
			wantErr: "failed to write key: open " +
				badKeyPath + ": no such file or directory",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			storage.keyPath = test.path

			err := storage.PutKey([]byte("key_data"))
			var gotErr string
			if err != nil {
				gotErr = err.Error()
			}
			assert.Equal(t, test.wantErr, gotErr, "Put key error")
		})
	}
}

func Test_Key(t *testing.T) {
	homeDir := testSetupEHomeDir(t)
	keyPath := filepath.Join(homeDir, parentDirName, storageDirName, keyFileName)
	badKeyPath := filepath.Join(homeDir, "not_exist", keyFileName)

	storage, err := NewFileSorage()
	require.Nil(t, err, "Create file storage")

	err = storage.PutKey([]byte("key_data"))
	require.Nil(t, err, "Set key value")

	tests := []struct {
		name    string
		path    string
		want    []byte
		wantErr string
	}{
		{name: "success", path: keyPath, want: []byte("key_data")},
		{
			name: "fail",
			path: badKeyPath,
			wantErr: "failed to get key: open " +
				badKeyPath + ": no such file or directory",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			storage.keyPath = test.path

			got, err := storage.Key()
			assert.Equal(t, test.want, got, "Get key")
			var gotErr string
			if err != nil {
				gotErr = err.Error()
			}
			assert.Equal(t, test.wantErr, gotErr, "Get key error")
		})
	}
}

func testSetupEHomeDir(t *testing.T) string {
	tmpDir := t.TempDir()
	// Linux / macOS (XDG)
	t.Setenv("XDG_CONFIG_HOME", tmpDir)
	// Fallback для macOS / Linux
	t.Setenv("HOME", tmpDir)
	// Windows
	t.Setenv("APPDATA", tmpDir)
	t.Setenv("USERPROFILE", tmpDir)

	return tmpDir
}
