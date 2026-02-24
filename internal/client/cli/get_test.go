package cli

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"github.com/EshkinKot1980/GophKeeper/internal/client/cli/mocks"
	"github.com/EshkinKot1980/GophKeeper/internal/common/dto"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_get(t *testing.T) {
	cr := dto.Credentials{Login: "user", Password: "password"}
	crData, err := json.Marshal(cr)
	require.Nil(t, err, "Credentials encode to json")

	type want struct {
		output string
		err    string
	}

	tests := []struct {
		name     string
		filePath string
		secretID string
		setup    func(t *testing.T) SecretService
		want     want
	}{
		{
			name:     "success_credentials",
			secretID: "13",
			setup: func(t *testing.T) SecretService {
				ctrl := gomock.NewController(t)
				service := mocks.NewMockSecretService(ctrl)
				service.EXPECT().
					GetSecretAndInfo(uint64(13)).
					Return(
						crData,
						dto.SecretInfo{
							Name:     "name",
							DataType: dto.SecretTypeCredentials,
						},
						nil,
					)
				return service
			},
			want: want{
				output: "name\n" +
					"--------------------------------\n" +
					"login:    user\n" +
					"password: password\n" +
					"--------------------------------\n" +
					"created: 0001-01-01 00:00:00\n",
			},
		},
		{
			name:     "credentials_invalid_json",
			secretID: "13",
			setup: func(t *testing.T) SecretService {
				ctrl := gomock.NewController(t)
				service := mocks.NewMockSecretService(ctrl)
				service.EXPECT().
					GetSecretAndInfo(uint64(13)).
					Return(
						[]byte("inalid json"),
						dto.SecretInfo{
							Name:     "name",
							DataType: dto.SecretTypeCredentials,
						},
						nil,
					)
				return service
			},
			want: want{
				err: "failed decode secret json: invalid character 'i' " +
					"looking for beginning of value",
			},
		},
		{
			name:     "success_text",
			secretID: "13",
			setup: func(t *testing.T) SecretService {
				ctrl := gomock.NewController(t)
				service := mocks.NewMockSecretService(ctrl)
				service.EXPECT().
					GetSecretAndInfo(uint64(13)).
					Return(
						[]byte("multiline\ntext"),
						dto.SecretInfo{
							Name:     "name",
							DataType: dto.SecretTypeText,
						},
						nil,
					)
				return service
			},
			want: want{
				output: "name\n" +
					"--------------------------------\n" +
					"multiline\n" +
					"text\n" +
					"--------------------------------\n" +
					"created: 0001-01-01 00:00:00\n",
			},
		},
		{
			name:     "file_invalid_path",
			secretID: "13",
			filePath: filepath.Join(t.TempDir(), "notexist", "file.bin"),
			setup: func(t *testing.T) SecretService {
				ctrl := gomock.NewController(t)
				service := mocks.NewMockSecretService(ctrl)
				service.EXPECT().
					GetSecretAndInfo(uint64(13)).
					Return(
						[]byte("file data"),
						dto.SecretInfo{
							Name:     "name",
							DataType: dto.SecretTypeFile,
						},
						nil,
					)
				return service
			},
			want: want{
				err: "failed to spot output directory",
			},
		},
		{
			name:     "unsuported_secret_type",
			secretID: "13",
			setup: func(t *testing.T) SecretService {
				ctrl := gomock.NewController(t)
				service := mocks.NewMockSecretService(ctrl)
				service.EXPECT().
					GetSecretAndInfo(uint64(13)).
					Return(
						[]byte("some data"),
						dto.SecretInfo{
							Name:     "name",
							DataType: "unknown",
						},
						nil,
					)
				return service
			},
			want: want{
				err: "unsuported secret type: unknown",
			},
		},
		{
			name:     "failed_to_get",
			secretID: "13",
			setup: func(t *testing.T) SecretService {
				ctrl := gomock.NewController(t)
				service := mocks.NewMockSecretService(ctrl)
				service.EXPECT().
					GetSecretAndInfo(uint64(13)).
					Return(nil, dto.SecretInfo{}, fmt.Errorf("failed to get secret"))
				return service
			},
			want: want{
				err: "failed to get secret",
			},
		},
		{
			name:     "invalid_id",
			secretID: "text_id",
			setup: func(t *testing.T) SecretService {
				ctrl := gomock.NewController(t)
				return mocks.NewMockSecretService(ctrl)
			},
			want: want{
				err: "id must be a number",
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			secretService = test.setup(t)
			filePath = test.filePath

			out := new(bytes.Buffer)
			err := get(out, test.secretID)

			var gotErr string
			if err != nil {
				gotErr = err.Error()
			}
			assert.Equal(t, test.want.err, gotErr, "Get error")
			assert.Equal(t, test.want.output, out.String(), "Get output")
		})
	}

}

func Test_saveFile(t *testing.T) {
	content := []byte("file content")

	tempDir := t.TempDir()
	err := os.Chdir(tempDir)
	require.Nil(t, err, "Change current dir")

	tests := []struct {
		name    string
		path    string
		info    dto.SecretInfo
		setup   func(t *testing.T) Prompt
		wantErr string
	}{
		{
			name: "success",
			path: filepath.Join(tempDir, "file.bin"),
			info: dto.SecretInfo{},
			setup: func(t *testing.T) Prompt {
				ctrl := gomock.NewController(t)
				return mocks.NewMockPrompt(ctrl)
			},
		},
		{
			name: "success_current_dir",
			info: dto.SecretInfo{Meta: []dto.MetaData{
				{Name: MetaFileName, Value: "meta.bin"},
			}},
			setup: func(t *testing.T) Prompt {
				ctrl := gomock.NewController(t)
				return mocks.NewMockPrompt(ctrl)
			},
		},
		{
			name: "success_current_dir_owerwrite",
			info: dto.SecretInfo{Meta: []dto.MetaData{
				{Name: MetaFileName, Value: "meta.bin"},
			}},
			setup: func(t *testing.T) Prompt {
				ctrl := gomock.NewController(t)
				prompt := mocks.NewMockPrompt(ctrl)
				prompt.EXPECT().
					Overwrite(gomock.All()).
					Return(true)
				return prompt
			},
		},
		{
			name: "success_current_dir_not_owerwrite",
			info: dto.SecretInfo{Meta: []dto.MetaData{
				{Name: MetaFileName, Value: "meta.bin"},
			}},
			setup: func(t *testing.T) Prompt {
				ctrl := gomock.NewController(t)
				prompt := mocks.NewMockPrompt(ctrl)
				prompt.EXPECT().
					Overwrite(gomock.All()).
					Return(false)
				return prompt
			},
		},
		{
			name: "success_path_not_owerwrite",
			path: filepath.Join(tempDir, "file.bin"),
			info: dto.SecretInfo{Meta: []dto.MetaData{}},
			setup: func(t *testing.T) Prompt {
				ctrl := gomock.NewController(t)
				prompt := mocks.NewMockPrompt(ctrl)
				prompt.EXPECT().
					Overwrite(gomock.All()).
					Return(false)
				return prompt
			},
		},
		{
			name: "invalid_path",
			path: "/some\x00dir/file.bin",
			info: dto.SecretInfo{},
			setup: func(t *testing.T) Prompt {
				ctrl := gomock.NewController(t)
				return mocks.NewMockPrompt(ctrl)
			},
			wantErr: "failed to get output file info: " +
				"stat /some\x00dir/file.bin: invalid argument",
		},
		{
			name: "dir_path_not_exist",
			path: filepath.Join("notexist", "file.bin"),
			info: dto.SecretInfo{},
			setup: func(t *testing.T) Prompt {
				ctrl := gomock.NewController(t)
				return mocks.NewMockPrompt(ctrl)
			},
			wantErr: "failed to spot output directory",
		},
		{
			name: "not_directory_in_path",
			path: filepath.Join("file.bin", "file.bin"),
			info: dto.SecretInfo{},
			setup: func(t *testing.T) Prompt {
				ctrl := gomock.NewController(t)
				return mocks.NewMockPrompt(ctrl)
			},
			wantErr: "failed to get output file info: " +
				"stat file.bin/file.bin: not a directory",
		},
		{
			name: "dir_without_metadata",
			path: tempDir,
			info: dto.SecretInfo{},
			setup: func(t *testing.T) Prompt {
				ctrl := gomock.NewController(t)
				return mocks.NewMockPrompt(ctrl)
			},
			wantErr: "failed to spot file name from metadata",
		},
		{
			name: "success_current_dir_not_owerwrite",
			info: dto.SecretInfo{Meta: []dto.MetaData{
				{Name: MetaFileName, Value: "meta.bin"},
			}},
			setup: func(t *testing.T) Prompt {
				ctrl := gomock.NewController(t)
				prompt := mocks.NewMockPrompt(ctrl)
				prompt.EXPECT().
					Overwrite(gomock.All()).
					Return(false)
				return prompt
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			filePath = test.path
			prompt = test.setup(t)

			err := saveFile(content, test.info)
			var gotErr string
			if err != nil {
				gotErr = err.Error()
			}
			assert.Equal(t, test.wantErr, gotErr, "Save file error")
		})
	}
}
