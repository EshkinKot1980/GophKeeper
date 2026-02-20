package service

import (
	"encoding/base64"
	"fmt"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/EshkinKot1980/GophKeeper/internal/client/service/mocks"
	"github.com/EshkinKot1980/GophKeeper/internal/common/crypto"
	"github.com/EshkinKot1980/GophKeeper/internal/common/dto"
)

func TestAuth_Register(t *testing.T) {
	token := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9" +
		".eyJleHAiOjE3NTg0NTk0OTMsImp0aSI6IjEifQ._mX-s6U9_iq4YhnQ5HOYbJAz7P8ly8BD_BufPYx2Kms"
	base64Salt := "EPhQ3C8pTRYDMac+aCzFTA"
	salt, err := base64.RawStdEncoding.DecodeString(base64Salt)
	require.Nil(t, err, "Decoding salt from base64")
	masterKey, err := crypto.DeriveKey([]byte("password13"), salt)
	require.Nil(t, err, "Generate masterKey from password")

	badBase64 := "EPhQ3C8pTRY&DMac+aCzFTA"

	tests := []struct {
		name    string
		cr      dto.Credentials
		cSetup  func(t *testing.T) Client
		sSetup  func(t *testing.T) Storage
		wantErr string
	}{
		{
			name: "success",
			cr:   dto.Credentials{Login: "test13", Password: "password13"},
			cSetup: func(t *testing.T) Client {
				ctrl := gomock.NewController(t)
				client := mocks.NewMockClient(ctrl)
				client.EXPECT().
					Register(dto.Credentials{Login: "test13", Password: "password13"}).
					Return(dto.AuthResponse{Token: token, EncrSalt: base64Salt}, nil)
				return client
			},
			sSetup: func(t *testing.T) Storage {
				ctrl := gomock.NewController(t)
				storage := mocks.NewMockStorage(ctrl)
				storage.EXPECT().PutKey(masterKey).Return(nil)
				storage.EXPECT().PutToken(token).Return(nil)
				return storage
			},
		},
		{
			name: "client_error",
			cr:   dto.Credentials{Login: "test13", Password: "password13"},
			cSetup: func(t *testing.T) Client {
				ctrl := gomock.NewController(t)
				client := mocks.NewMockClient(ctrl)
				client.EXPECT().
					Register(dto.Credentials{Login: "test13", Password: "password13"}).
					Return(dto.AuthResponse{}, fmt.Errorf("registration failed"))
				return client
			},
			sSetup: func(t *testing.T) Storage {
				ctrl := gomock.NewController(t)
				return mocks.NewMockStorage(ctrl)
			},
			wantErr: "registration failed",
		},
		{
			name: "invalid_salt",
			cr:   dto.Credentials{Login: "test13", Password: "password13"},
			cSetup: func(t *testing.T) Client {
				ctrl := gomock.NewController(t)
				client := mocks.NewMockClient(ctrl)
				client.EXPECT().
					Register(dto.Credentials{Login: "test13", Password: "password13"}).
					Return(dto.AuthResponse{Token: token, EncrSalt: badBase64}, nil)
				return client
			},
			sSetup: func(t *testing.T) Storage {
				ctrl := gomock.NewController(t)
				return mocks.NewMockStorage(ctrl)
			},
			wantErr: "the server returned invalid data",
		},
		{
			name: "derive_key_error",
			cr:   dto.Credentials{Login: "test13", Password: "password13"},
			cSetup: func(t *testing.T) Client {
				ctrl := gomock.NewController(t)
				client := mocks.NewMockClient(ctrl)
				client.EXPECT().
					Register(dto.Credentials{Login: "test13", Password: "password13"}).
					Return(dto.AuthResponse{Token: token, EncrSalt: ""}, nil)
				return client
			},
			sSetup: func(t *testing.T) Storage {
				ctrl := gomock.NewController(t)
				return mocks.NewMockStorage(ctrl)
			},
			wantErr: "failed to derive key",
		},
		{
			name: "store_key_error",
			cr:   dto.Credentials{Login: "test13", Password: "password13"},
			cSetup: func(t *testing.T) Client {
				ctrl := gomock.NewController(t)
				client := mocks.NewMockClient(ctrl)
				client.EXPECT().
					Register(dto.Credentials{Login: "test13", Password: "password13"}).
					Return(dto.AuthResponse{Token: token, EncrSalt: base64Salt}, nil)
				return client
			},
			sSetup: func(t *testing.T) Storage {
				ctrl := gomock.NewController(t)
				storage := mocks.NewMockStorage(ctrl)
				storage.EXPECT().
					PutKey(masterKey).Return(fmt.Errorf("failed to put key"))
				return storage
			},
			wantErr: "failed to store key",
		},
		{
			name: "store_token_error",
			cr:   dto.Credentials{Login: "test13", Password: "password13"},
			cSetup: func(t *testing.T) Client {
				ctrl := gomock.NewController(t)
				client := mocks.NewMockClient(ctrl)
				client.EXPECT().
					Register(dto.Credentials{Login: "test13", Password: "password13"}).
					Return(dto.AuthResponse{Token: token, EncrSalt: base64Salt}, nil)
				return client
			},
			sSetup: func(t *testing.T) Storage {
				ctrl := gomock.NewController(t)
				storage := mocks.NewMockStorage(ctrl)
				storage.EXPECT().PutKey(masterKey).Return(nil)
				storage.EXPECT().
					PutToken(token).Return(fmt.Errorf("failed to put token"))
				return storage
			},
			wantErr: "failed to store token",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			client := test.cSetup(t)
			storage := test.sSetup(t)

			authService := NewAuth(client, storage)
			err := authService.Register(test.cr)

			var gotErr string
			if err != nil {
				gotErr = err.Error()
			}

			assert.Equal(t, test.wantErr, gotErr, "Register error")
		})
	}
}

func TestAuth_Login(t *testing.T) {
	token := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9" +
		".eyJleHAiOjE3NTg0NTk0OTMsImp0aSI6IjEifQ._mX-s6U9_iq4YhnQ5HOYbJAz7P8ly8BD_BufPYx2Kms"
	base64Salt := "EPhQ3C8pTRYDMac+aCzFTA"
	salt, err := base64.RawStdEncoding.DecodeString(base64Salt)
	require.Nil(t, err, "Decoding salt from base64")
	masterKey, err := crypto.DeriveKey([]byte("password13"), salt)
	require.Nil(t, err, "Generate masterKey from password")

	tests := []struct {
		name    string
		cr      dto.Credentials
		cSetup  func(t *testing.T) Client
		sSetup  func(t *testing.T) Storage
		wantErr string
	}{
		{
			name: "success",
			cr:   dto.Credentials{Login: "test13", Password: "password13"},
			cSetup: func(t *testing.T) Client {
				ctrl := gomock.NewController(t)
				client := mocks.NewMockClient(ctrl)
				client.EXPECT().
					Login(dto.Credentials{Login: "test13", Password: "password13"}).
					Return(dto.AuthResponse{Token: token, EncrSalt: base64Salt}, nil)
				return client
			},
			sSetup: func(t *testing.T) Storage {
				ctrl := gomock.NewController(t)
				storage := mocks.NewMockStorage(ctrl)
				storage.EXPECT().PutKey(masterKey).Return(nil)
				storage.EXPECT().PutToken(token).Return(nil)
				return storage
			},
		},
		{
			name: "client_error",
			cr:   dto.Credentials{Login: "test13", Password: "password13"},
			cSetup: func(t *testing.T) Client {
				ctrl := gomock.NewController(t)
				client := mocks.NewMockClient(ctrl)
				client.EXPECT().
					Login(dto.Credentials{Login: "test13", Password: "password13"}).
					Return(dto.AuthResponse{}, fmt.Errorf("login failed"))
				return client
			},
			sSetup: func(t *testing.T) Storage {
				ctrl := gomock.NewController(t)
				return mocks.NewMockStorage(ctrl)
			},
			wantErr: "login failed",
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			client := test.cSetup(t)
			storage := test.sSetup(t)

			authService := NewAuth(client, storage)
			err := authService.Login(test.cr)

			var gotErr string
			if err != nil {
				gotErr = err.Error()
			}

			assert.Equal(t, test.wantErr, gotErr, "Register error")
		})
	}
}
