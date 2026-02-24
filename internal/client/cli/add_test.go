package cli

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"github.com/EshkinKot1980/GophKeeper/internal/client/cli/mocks"
	"github.com/EshkinKot1980/GophKeeper/internal/client/config"
	"github.com/EshkinKot1980/GophKeeper/internal/common/dto"
	"github.com/golang/mock/gomock"
	"github.com/spf13/cobra"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_addCredentials(t *testing.T) {
	credentials := dto.Credentials{Login: "login", Password: "password1"}
	data, err := json.Marshal(credentials)
	require.Nil(t, err, "Credentials json encode")

	tests := []struct {
		name    string
		pSetup  func(t *testing.T) Prompt
		sSetup  func(t *testing.T) SecretService
		wantErr string
	}{
		{
			name: "success",
			pSetup: func(t *testing.T) Prompt {
				ctrl := gomock.NewController(t)
				prompt := mocks.NewMockPrompt(ctrl)
				prompt.EXPECT().
					SecretName().Return("name", nil)
				prompt.EXPECT().
					Credentials().Return(credentials, nil)
				return prompt
			},
			sSetup: func(t *testing.T) SecretService {
				ctrl := gomock.NewController(t)
				service := mocks.NewMockSecretService(ctrl)
				service.EXPECT().
					Upload(
						dto.SecretRequest{
							Name:     "name",
							DataType: dto.SecretTypeCredentials,
							Meta:     []dto.MetaData{},
						},
						data,
					).
					Return(nil)
				return service
			},
		},
		{
			name: "invalid_name",
			pSetup: func(t *testing.T) Prompt {
				ctrl := gomock.NewController(t)
				prompt := mocks.NewMockPrompt(ctrl)
				prompt.EXPECT().
					SecretName().Return("", fmt.Errorf("invalid name"))
				return prompt
			},
			sSetup: func(t *testing.T) SecretService {
				ctrl := gomock.NewController(t)
				return mocks.NewMockSecretService(ctrl)
			},
			wantErr: "invalid name",
		},
		{
			name: "invalid_credentials",
			pSetup: func(t *testing.T) Prompt {
				ctrl := gomock.NewController(t)
				prompt := mocks.NewMockPrompt(ctrl)
				prompt.EXPECT().
					SecretName().Return("name", nil)
				prompt.EXPECT().
					Credentials().Return(dto.Credentials{}, fmt.Errorf("invalid credentials"))
				return prompt
			},
			sSetup: func(t *testing.T) SecretService {
				ctrl := gomock.NewController(t)
				return mocks.NewMockSecretService(ctrl)
			},
			wantErr: "invalid credentials",
		},
		{
			name: "failed_to_send",
			pSetup: func(t *testing.T) Prompt {
				ctrl := gomock.NewController(t)
				prompt := mocks.NewMockPrompt(ctrl)
				prompt.EXPECT().
					SecretName().Return("name", nil)
				prompt.EXPECT().
					Credentials().Return(credentials, nil)
				return prompt
			},
			sSetup: func(t *testing.T) SecretService {
				ctrl := gomock.NewController(t)
				service := mocks.NewMockSecretService(ctrl)
				service.EXPECT().
					Upload(gomock.All(), gomock.All()).
					Return(fmt.Errorf("authorization failed"))
				return service
			},
			wantErr: "failed to send data to server: authorization failed",
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			prompt = test.pSetup(t)
			secretService = test.sSetup(t)

			err := addCredentials(&cobra.Command{}, []string{})
			var gotErr string
			if err != nil {
				gotErr = err.Error()
			}
			assert.Equal(t, test.wantErr, gotErr, "Add clredentials error")
		})
	}
}

func Test_addFile(t *testing.T) {
	tmpdir := t.TempDir()

	dirPath := filepath.Join(tmpdir, "dirName")
	err := os.MkdirAll(dirPath, 0755)
	require.Nil(t, err, "Dir creating")

	goodPath := filepath.Join(tmpdir, "file.bin")
	err = os.WriteFile(goodPath, []byte("file data"), 0600)
	require.Nil(t, err, "File creating")

	meta := []dto.MetaData{
		{Name: MetaFileName, Value: "file.bin"},
		{Name: MetaFilePath, Value: tmpdir},
	}

	tests := []struct {
		name    string
		path    string
		maxSize int64
		pSetup  func(t *testing.T) Prompt
		sSetup  func(t *testing.T) SecretService
		wantErr string
	}{
		{
			name:    "success",
			path:    goodPath,
			maxSize: 4096,
			pSetup: func(t *testing.T) Prompt {
				ctrl := gomock.NewController(t)
				prompt := mocks.NewMockPrompt(ctrl)
				prompt.EXPECT().
					SecretName().Return("name", nil)
				return prompt
			},
			sSetup: func(t *testing.T) SecretService {
				ctrl := gomock.NewController(t)
				service := mocks.NewMockSecretService(ctrl)
				service.EXPECT().
					Upload(
						dto.SecretRequest{
							Name:     "name",
							DataType: dto.SecretTypeFile,
							Meta:     meta,
						},
						[]byte("file data"),
					).
					Return(nil)
				return service
			},
		},
		{
			name:    "not_exist",
			path:    filepath.Join(tmpdir, "notExist.txt"),
			maxSize: 4096,
			pSetup: func(t *testing.T) Prompt {
				ctrl := gomock.NewController(t)
				return mocks.NewMockPrompt(ctrl)
			},
			sSetup: func(t *testing.T) SecretService {
				ctrl := gomock.NewController(t)
				return mocks.NewMockSecretService(ctrl)
			},
			wantErr: fmt.Sprintf(
				"failed to get file info: stat %s: no such file or directory",
				filepath.Join(tmpdir, "notExist.txt"),
			),
		},
		{
			name:    "is_directory",
			path:    dirPath,
			maxSize: 4096,
			pSetup: func(t *testing.T) Prompt {
				ctrl := gomock.NewController(t)
				return mocks.NewMockPrompt(ctrl)
			},
			sSetup: func(t *testing.T) SecretService {
				ctrl := gomock.NewController(t)
				return mocks.NewMockSecretService(ctrl)
			},
			wantErr: "failed add file: the file \"dirName\" is directory",
		},
		{
			name:    "too_large",
			path:    goodPath,
			maxSize: 4,
			pSetup: func(t *testing.T) Prompt {
				ctrl := gomock.NewController(t)
				return mocks.NewMockPrompt(ctrl)
			},
			sSetup: func(t *testing.T) SecretService {
				ctrl := gomock.NewController(t)
				return mocks.NewMockSecretService(ctrl)
			},
			wantErr: "file size (9 bytes) exceeds the limit of 4 bytes",
		},
		{
			name:    "invalid_secret_name",
			path:    goodPath,
			maxSize: 4096,
			pSetup: func(t *testing.T) Prompt {
				ctrl := gomock.NewController(t)
				prompt := mocks.NewMockPrompt(ctrl)
				prompt.EXPECT().
					SecretName().Return("", fmt.Errorf("invalid name"))
				return prompt
			},
			sSetup: func(t *testing.T) SecretService {
				ctrl := gomock.NewController(t)
				return mocks.NewMockSecretService(ctrl)
			},
			wantErr: "invalid name",
		},
		{
			name:    "failed_to_send",
			path:    goodPath,
			maxSize: 4096,
			pSetup: func(t *testing.T) Prompt {
				ctrl := gomock.NewController(t)
				prompt := mocks.NewMockPrompt(ctrl)
				prompt.EXPECT().
					SecretName().Return("name", nil)
				return prompt
			},
			sSetup: func(t *testing.T) SecretService {
				ctrl := gomock.NewController(t)
				service := mocks.NewMockSecretService(ctrl)
				service.EXPECT().
					Upload(gomock.All(), gomock.All()).
					Return(fmt.Errorf("sending error"))
				return service
			},
			wantErr: "failed to send data to server: sending error",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			prompt = test.pSetup(t)
			secretService = test.sSetup(t)
			cfg = &config.Config{FileMaxSize: test.maxSize}

			out := new(bytes.Buffer)
			err := addFile(out, test.path)

			var gotErr string
			if err != nil {
				gotErr = err.Error()
			}
			assert.Equal(t, test.wantErr, gotErr, "Add file error")
		})
	}
}

func Test_addText(t *testing.T) {
	text := `
	some multiline text
	some multiline text
	`

	tests := []struct {
		name    string
		pSetup  func(t *testing.T) Prompt
		sSetup  func(t *testing.T) SecretService
		wantErr string
	}{
		{
			name: "success",
			pSetup: func(t *testing.T) Prompt {
				ctrl := gomock.NewController(t)
				prompt := mocks.NewMockPrompt(ctrl)
				prompt.EXPECT().
					SecretName().Return("name", nil)
				prompt.EXPECT().
					Text().Return(text, nil)
				return prompt
			},
			sSetup: func(t *testing.T) SecretService {
				ctrl := gomock.NewController(t)
				service := mocks.NewMockSecretService(ctrl)
				service.EXPECT().
					Upload(
						dto.SecretRequest{
							Name:     "name",
							DataType: dto.SecretTypeText,
							Meta:     []dto.MetaData{},
						},
						[]byte(text),
					).
					Return(nil)
				return service
			},
		},
		{
			name: "invalid_name",
			pSetup: func(t *testing.T) Prompt {
				ctrl := gomock.NewController(t)
				prompt := mocks.NewMockPrompt(ctrl)
				prompt.EXPECT().
					SecretName().Return("", fmt.Errorf("invalid name"))
				return prompt
			},
			sSetup: func(t *testing.T) SecretService {
				ctrl := gomock.NewController(t)
				return mocks.NewMockSecretService(ctrl)
			},
			wantErr: "invalid name",
		},
		{
			name: "text_scan_error",
			pSetup: func(t *testing.T) Prompt {
				ctrl := gomock.NewController(t)
				prompt := mocks.NewMockPrompt(ctrl)
				prompt.EXPECT().
					SecretName().Return("name", nil)
				prompt.EXPECT().
					Text().Return("", fmt.Errorf("text scan error"))
				return prompt
			},
			sSetup: func(t *testing.T) SecretService {
				ctrl := gomock.NewController(t)
				return mocks.NewMockSecretService(ctrl)
			},
			wantErr: "text scan error",
		},
		{
			name: "failed_to_send",
			pSetup: func(t *testing.T) Prompt {
				ctrl := gomock.NewController(t)
				prompt := mocks.NewMockPrompt(ctrl)
				prompt.EXPECT().
					SecretName().Return("name", nil)
				prompt.EXPECT().
					Text().Return(text, nil)
				return prompt
			},
			sSetup: func(t *testing.T) SecretService {
				ctrl := gomock.NewController(t)
				service := mocks.NewMockSecretService(ctrl)
				service.EXPECT().
					Upload(gomock.All(), gomock.All()).
					Return(fmt.Errorf("authorization failed"))
				return service
			},
			wantErr: "failed to send data to server: authorization failed",
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			prompt = test.pSetup(t)
			secretService = test.sSetup(t)

			err := addText(&cobra.Command{}, []string{})
			var gotErr string
			if err != nil {
				gotErr = err.Error()
			}
			assert.Equal(t, test.wantErr, gotErr, "Add text error")
		})
	}
}
