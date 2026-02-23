package service

import (
	"encoding/base64"
	"fmt"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	httpClient "github.com/EshkinKot1980/GophKeeper/internal/client/http"
	"github.com/EshkinKot1980/GophKeeper/internal/client/service/mocks"
	"github.com/EshkinKot1980/GophKeeper/internal/common/crypto"
	"github.com/EshkinKot1980/GophKeeper/internal/common/dto"
)

func TestSecret_Upload(t *testing.T) {
	rawSecret := dto.SecretRequest{Name: "test", DataType: dto.SecretTypeText}
	token := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9" +
		".eyJleHAiOjE3NTg0NTk0OTMsImp0aSI6IjEifQ._mX-s6U9_iq4YhnQ5HOYbJAz7P8ly8BD_BufPYx2Kms"
	base64Salt := "EPhQ3C8pTRYDMac+aCzFTA"
	salt, err := base64.RawStdEncoding.DecodeString(base64Salt)
	require.Nil(t, err, "Decoding salt from base64")
	masterKey, err := crypto.DeriveKey([]byte("password13"), salt)
	require.Nil(t, err, "Master key creation")

	tests := []struct {
		name    string
		dto     dto.SecretRequest
		data    []byte
		sSetup  func(t *testing.T) Storage
		cSetup  func(t *testing.T) Client
		wantErr string
	}{
		{
			name: "success",
			dto:  rawSecret,
			data: []byte("data to crypt"),
			sSetup: func(t *testing.T) Storage {
				ctrl := gomock.NewController(t)
				storage := mocks.NewMockStorage(ctrl)
				storage.EXPECT().
					Key().Return(masterKey, nil)
				storage.EXPECT().
					Token().Return(token, nil)
				return storage
			},
			cSetup: func(t *testing.T) Client {
				ctrl := gomock.NewController(t)
				client := mocks.NewMockClient(ctrl)
				client.EXPECT().
					Upload(gomock.All(), token).
					Return(nil)
				return client
			},
		},
		{
			name: "without_master_key",
			dto:  rawSecret,
			data: []byte("data to crypt"),
			sSetup: func(t *testing.T) Storage {
				ctrl := gomock.NewController(t)
				storage := mocks.NewMockStorage(ctrl)
				storage.EXPECT().
					Key().Return(nil, fmt.Errorf("any error"))
				return storage
			},
			cSetup: func(t *testing.T) Client {
				ctrl := gomock.NewController(t)
				return mocks.NewMockClient(ctrl)
			},
			wantErr: "authorization failed: any error",
		},
		{
			name: "bad_master_key",
			dto:  rawSecret,
			data: []byte("data to crypt"),
			sSetup: func(t *testing.T) Storage {
				ctrl := gomock.NewController(t)
				storage := mocks.NewMockStorage(ctrl)
				storage.EXPECT().
					Key().Return([]byte{}, nil)
				return storage
			},
			cSetup: func(t *testing.T) Client {
				ctrl := gomock.NewController(t)
				return mocks.NewMockClient(ctrl)
			},
			wantErr: "failed to encrypt secret: failed to encrypt DEK:" +
				" failed to create cipher block: crypto/aes: invalid key size 0",
		},
		{
			name: "without_token",
			dto:  rawSecret,
			data: []byte("data to crypt"),
			sSetup: func(t *testing.T) Storage {
				ctrl := gomock.NewController(t)
				storage := mocks.NewMockStorage(ctrl)
				storage.EXPECT().
					Key().Return(masterKey, nil)
				storage.EXPECT().
					Token().Return("", fmt.Errorf("any error"))
				return storage
			},
			cSetup: func(t *testing.T) Client {
				ctrl := gomock.NewController(t)
				return mocks.NewMockClient(ctrl)
			},
			wantErr: "authorization failed: any error",
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			client := test.cSetup(t)
			storage := test.sSetup(t)

			secretService := NewSecret(client, storage)
			err := secretService.Upload(test.dto, test.data)

			var gotErr string
			if err != nil {
				gotErr = err.Error()
			}
			assert.Equal(t, test.wantErr, gotErr, "Upload error")
		})
	}
}

func TestSecret_GetSecretAndInfo(t *testing.T) {
	token := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9" +
		".eyJleHAiOjE3NTg0NTk0OTMsImp0aSI6IjEifQ._mX-s6U9_iq4YhnQ5HOYbJAz7P8ly8BD_BufPYx2Kms"

	base64Salt := "EPhQ3C8pTRYDMac+aCzFTA"
	salt, err := base64.RawStdEncoding.DecodeString(base64Salt)
	require.Nil(t, err, "Decoding salt from base64")
	masterKey, err := crypto.DeriveKey([]byte("password13"), salt)
	require.Nil(t, err, "Master key creation")

	secret := []byte("some secret text")
	ecnrData, err := encryptData(masterKey, secret)
	require.Nil(t, err, "Data ecryption")

	type want struct {
		secret []byte
		info   dto.SecretInfo
		err    error
	}

	tests := []struct {
		name     string
		secretID uint64
		sSetup   func(t *testing.T) Storage
		cSetup   func(t *testing.T) Client
		want     want
	}{
		{
			name:     "success",
			secretID: 13,
			sSetup: func(t *testing.T) Storage {
				ctrl := gomock.NewController(t)
				storage := mocks.NewMockStorage(ctrl)
				storage.EXPECT().
					Key().Return(masterKey, nil)
				storage.EXPECT().
					Token().Return(token, nil)
				return storage
			},
			cSetup: func(t *testing.T) Client {
				ctrl := gomock.NewController(t)
				client := mocks.NewMockClient(ctrl)
				client.EXPECT().
					Retrieve(uint64(13), token).
					Return(dto.SecretResponse{EncrData: ecnrData}, nil)
				return client
			},
			want: want{
				secret: secret,
				info:   dto.SecretInfo{},
			},
		},
		{
			name:     "without_token",
			secretID: 13,
			sSetup: func(t *testing.T) Storage {
				ctrl := gomock.NewController(t)
				storage := mocks.NewMockStorage(ctrl)
				storage.EXPECT().
					Key().Return(masterKey, nil)
				storage.EXPECT().
					Token().Return("", fmt.Errorf("any error"))
				return storage
			},
			cSetup: func(t *testing.T) Client {
				ctrl := gomock.NewController(t)
				return mocks.NewMockClient(ctrl)
			},
			want: want{
				err: ErrAuthorizationFailed,
			},
		},
		{
			name:     "without_master_key",
			secretID: 13,
			sSetup: func(t *testing.T) Storage {
				ctrl := gomock.NewController(t)
				storage := mocks.NewMockStorage(ctrl)
				storage.EXPECT().
					Key().Return(nil, fmt.Errorf("any error"))
				return storage
			},
			cSetup: func(t *testing.T) Client {
				ctrl := gomock.NewController(t)
				return mocks.NewMockClient(ctrl)
			},
			want: want{
				err: ErrAuthorizationFailed,
			},
		},
		{
			name:     "bad_master_key",
			secretID: 13,
			sSetup: func(t *testing.T) Storage {
				ctrl := gomock.NewController(t)
				storage := mocks.NewMockStorage(ctrl)
				storage.EXPECT().
					Key().Return([]byte{}, nil)
				storage.EXPECT().
					Token().Return(token, nil)
				return storage
			},
			cSetup: func(t *testing.T) Client {
				ctrl := gomock.NewController(t)
				client := mocks.NewMockClient(ctrl)
				client.EXPECT().
					Retrieve(uint64(13), token).
					Return(dto.SecretResponse{EncrData: ecnrData}, nil)
				return client
			},
			want: want{
				err: ErrSecretDecryptionFailed,
			},
		},
		{
			name:     "client_err",
			secretID: 13,
			sSetup: func(t *testing.T) Storage {
				ctrl := gomock.NewController(t)
				storage := mocks.NewMockStorage(ctrl)
				storage.EXPECT().
					Key().Return([]byte{}, nil)
				storage.EXPECT().
					Token().Return(token, nil)
				return storage
			},
			cSetup: func(t *testing.T) Client {
				ctrl := gomock.NewController(t)
				client := mocks.NewMockClient(ctrl)
				client.EXPECT().
					Retrieve(uint64(13), token).
					Return(dto.SecretResponse{}, httpClient.ErrSecretRetrieveFailed)
				return client
			},
			want: want{
				err: httpClient.ErrSecretRetrieveFailed,
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			client := test.cSetup(t)
			storage := test.sSetup(t)

			secretService := NewSecret(client, storage)
			secret, info, err := secretService.GetSecretAndInfo(test.secretID)

			assert.ErrorIs(t, err, test.want.err, "Get secret and info error")
			if err != nil {
				return
			}
			assert.Equal(t, test.want.secret, secret, "Get secret")
			assert.Equal(t, test.want.info, info, "Get secret info")
		})
	}
}

func TestSecret_InfoList(t *testing.T) {
	token := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9" +
		".eyJleHAiOjE3NTg0NTk0OTMsImp0aSI6IjEifQ._mX-s6U9_iq4YhnQ5HOYbJAz7P8ly8BD_BufPYx2Kms"
	infoList := []dto.SecretInfo{{}}

	type want struct {
		list []dto.SecretInfo
		err  error
	}

	tests := []struct {
		name   string
		sSetup func(t *testing.T) Storage
		cSetup func(t *testing.T) Client
		want   want
	}{
		{
			name: "success",
			sSetup: func(t *testing.T) Storage {
				ctrl := gomock.NewController(t)
				storage := mocks.NewMockStorage(ctrl)
				storage.EXPECT().
					Token().Return(token, nil)
				return storage
			},
			cSetup: func(t *testing.T) Client {
				ctrl := gomock.NewController(t)
				client := mocks.NewMockClient(ctrl)
				client.EXPECT().
					InfoList(token).
					Return(infoList, nil)
				return client
			},
			want: want{
				list: infoList,
			},
		},
		{
			name: "without_token",
			sSetup: func(t *testing.T) Storage {
				ctrl := gomock.NewController(t)
				storage := mocks.NewMockStorage(ctrl)
				storage.EXPECT().
					Token().Return("", fmt.Errorf("any error"))
				return storage
			},
			cSetup: func(t *testing.T) Client {
				ctrl := gomock.NewController(t)
				return mocks.NewMockClient(ctrl)
			},
			want: want{
				err: ErrAuthorizationFailed,
			},
		},

		{
			name: "client_err",
			sSetup: func(t *testing.T) Storage {
				ctrl := gomock.NewController(t)
				storage := mocks.NewMockStorage(ctrl)
				storage.EXPECT().
					Token().Return(token, nil)
				return storage
			},
			cSetup: func(t *testing.T) Client {
				ctrl := gomock.NewController(t)
				client := mocks.NewMockClient(ctrl)
				client.EXPECT().
					InfoList(token).
					Return(nil, httpClient.ErrSecretInfoListFailed)
				return client
			},
			want: want{
				err: httpClient.ErrSecretInfoListFailed,
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			client := test.cSetup(t)
			storage := test.sSetup(t)

			secretService := NewSecret(client, storage)
			list, err := secretService.InfoList()

			assert.ErrorIs(t, err, test.want.err, "Get secret and info error")
			if err != nil {
				return
			}
			assert.Equal(t, test.want.list, list, "Get secret list")
		})
	}
}
