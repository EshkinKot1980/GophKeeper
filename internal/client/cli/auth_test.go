package cli

import (
	"fmt"
	"testing"

	"github.com/EshkinKot1980/GophKeeper/internal/client/cli/mocks"
	"github.com/EshkinKot1980/GophKeeper/internal/common/dto"
	"github.com/golang/mock/gomock"
	"github.com/spf13/cobra"
	"github.com/stretchr/testify/assert"
)

func Test_login(t *testing.T) {
	tests := []struct {
		name    string
		pSetup  func(t *testing.T) Prompt
		sSetup  func(t *testing.T) AuthService
		wantErr string
	}{
		{
			name: "success",
			pSetup: func(t *testing.T) Prompt {
				ctrl := gomock.NewController(t)
				prompt := mocks.NewMockPrompt(ctrl)
				prompt.EXPECT().
					Credentials().Return(dto.Credentials{}, nil)
				return prompt
			},
			sSetup: func(t *testing.T) AuthService {
				ctrl := gomock.NewController(t)
				service := mocks.NewMockAuthService(ctrl)
				service.EXPECT().
					Login(gomock.All()).Return(nil)
				return service
			},
		},
		{
			name: "invalid_credentials",
			pSetup: func(t *testing.T) Prompt {
				ctrl := gomock.NewController(t)
				prompt := mocks.NewMockPrompt(ctrl)
				prompt.EXPECT().
					Credentials().
					Return(dto.Credentials{}, fmt.Errorf("some err"))
				return prompt
			},
			sSetup: func(t *testing.T) AuthService {
				ctrl := gomock.NewController(t)
				return mocks.NewMockAuthService(ctrl)
			},
			wantErr: "some err",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			prompt = test.pSetup(t)
			authService = test.sSetup(t)

			err := login(&cobra.Command{}, []string{})
			var gotErr string
			if err != nil {
				gotErr = err.Error()
			}
			assert.Equal(t, test.wantErr, gotErr, "Login error")
		})
	}
}

func Test_Register(t *testing.T) {
	tests := []struct {
		name    string
		pSetup  func(t *testing.T) Prompt
		sSetup  func(t *testing.T) AuthService
		wantErr string
	}{
		{
			name: "success",
			pSetup: func(t *testing.T) Prompt {
				ctrl := gomock.NewController(t)
				prompt := mocks.NewMockPrompt(ctrl)
				prompt.EXPECT().
					RegisterCredentials().
					Return(dto.Credentials{}, nil)
				return prompt
			},
			sSetup: func(t *testing.T) AuthService {
				ctrl := gomock.NewController(t)
				service := mocks.NewMockAuthService(ctrl)
				service.EXPECT().
					Register(gomock.All()).Return(nil)
				return service
			},
		},
		{
			name: "invalid_credentials",
			pSetup: func(t *testing.T) Prompt {
				ctrl := gomock.NewController(t)
				prompt := mocks.NewMockPrompt(ctrl)
				prompt.EXPECT().
					RegisterCredentials().
					Return(dto.Credentials{}, fmt.Errorf("some err"))
				return prompt
			},
			sSetup: func(t *testing.T) AuthService {
				ctrl := gomock.NewController(t)
				return mocks.NewMockAuthService(ctrl)
			},
			wantErr: "some err",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			prompt = test.pSetup(t)
			authService = test.sSetup(t)

			err := register(&cobra.Command{}, []string{})
			var gotErr string
			if err != nil {
				gotErr = err.Error()
			}
			assert.Equal(t, test.wantErr, gotErr, "Register error")
		})
	}
}
