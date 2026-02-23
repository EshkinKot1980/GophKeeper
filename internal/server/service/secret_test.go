package service

import (
	"context"
	"fmt"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"

	"github.com/EshkinKot1980/GophKeeper/internal/common/dto"
	"github.com/EshkinKot1980/GophKeeper/internal/server/entity"

	repErrors "github.com/EshkinKot1980/GophKeeper/internal/server/repository/errors"
	srvContext "github.com/EshkinKot1980/GophKeeper/internal/server/service/context"
	srvErrors "github.com/EshkinKot1980/GophKeeper/internal/server/service/errors"
	"github.com/EshkinKot1980/GophKeeper/internal/server/service/mocks"
)

func TestSecret_Save(t *testing.T) {
	userID := "1ed655b6-0738-4162-a34a-34257c0dc106"
	goodCtx := srvContext.SetUserID(context.Background(), userID)
	requestDTO := dto.SecretRequest{}

	tests := []struct {
		name    string
		ctx     context.Context
		secret  *dto.SecretRequest
		rSetup  func(t *testing.T) SecretRepository
		lSetup  func(t *testing.T) Logger
		wantErr error
	}{
		{
			name:   "success",
			ctx:    goodCtx,
			secret: &requestDTO,
			rSetup: func(t *testing.T) SecretRepository {
				ctrl := gomock.NewController(t)
				repository := mocks.NewMockSecretRepository(ctrl)
				repository.EXPECT().
					Create(gomock.All(), gomock.All()).
					Return(nil)
				return repository
			},
			lSetup: func(t *testing.T) Logger {
				ctrl := gomock.NewController(t)
				return mocks.NewMockLogger(ctrl)
			},
		},
		{
			name:   "witout_user",
			ctx:    context.TODO(),
			secret: &requestDTO,
			rSetup: func(t *testing.T) SecretRepository {
				ctrl := gomock.NewController(t)
				return mocks.NewMockSecretRepository(ctrl)
			},
			lSetup: func(t *testing.T) Logger {
				ctrl := gomock.NewController(t)
				logger := mocks.NewMockLogger(ctrl)
				logger.EXPECT().
					Error("failed to get user id", gomock.All())
				return logger
			},
			wantErr: srvErrors.ErrUnexpected,
		},
		{
			name:   "repository_error",
			ctx:    goodCtx,
			secret: &requestDTO,
			rSetup: func(t *testing.T) SecretRepository {
				ctrl := gomock.NewController(t)
				repository := mocks.NewMockSecretRepository(ctrl)
				repository.EXPECT().
					Create(gomock.All(), gomock.All()).
					Return(fmt.Errorf("repositoryerror"))
				return repository
			},
			lSetup: func(t *testing.T) Logger {
				ctrl := gomock.NewController(t)
				logger := mocks.NewMockLogger(ctrl)
				logger.EXPECT().
					Error("failed create secret", gomock.All())
				return logger
			},
			wantErr: srvErrors.ErrUnexpected,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			repository := test.rSetup(t)
			logger := test.lSetup(t)
			secretService := NewSecret(logger, repository)
			err := secretService.Save(test.ctx, test.secret)
			assert.ErrorIs(t, err, test.wantErr, "Save secret error")
		})
	}
}

func TestSecret_Secret(t *testing.T) {
	userID := "1ed655b6-0738-4162-a34a-34257c0dc106"
	goodCtx := srvContext.SetUserID(context.Background(), userID)

	type want struct {
		secret dto.SecretResponse
		err    error
	}

	tests := []struct {
		name     string
		ctx      context.Context
		secretID uint64
		rSetup   func(t *testing.T) SecretRepository
		lSetup   func(t *testing.T) Logger
		want     want
	}{
		{
			name:     "success",
			ctx:      goodCtx,
			secretID: 13,
			rSetup: func(t *testing.T) SecretRepository {
				ctrl := gomock.NewController(t)
				repository := mocks.NewMockSecretRepository(ctrl)
				repository.EXPECT().
					GetForUser(gomock.All(), gomock.All(), gomock.All()).
					Return(entity.Secret{MetaData: "[]"}, nil)
				return repository
			},
			lSetup: func(t *testing.T) Logger {
				ctrl := gomock.NewController(t)
				return mocks.NewMockLogger(ctrl)
			},
			want: want{
				secret: dto.SecretResponse{Meta: []dto.MetaData{}},
			},
		},
		{
			name:     "without_user",
			ctx:      context.TODO(),
			secretID: 13,
			rSetup: func(t *testing.T) SecretRepository {
				ctrl := gomock.NewController(t)
				return mocks.NewMockSecretRepository(ctrl)
			},
			lSetup: func(t *testing.T) Logger {
				ctrl := gomock.NewController(t)
				logger := mocks.NewMockLogger(ctrl)
				logger.EXPECT().
					Error("failed to get user id", gomock.All())
				return logger
			},
			want: want{
				err: srvErrors.ErrUnexpected,
			},
		},
		{
			name:     "bad_meta_data",
			ctx:      goodCtx,
			secretID: 13,
			rSetup: func(t *testing.T) SecretRepository {
				ctrl := gomock.NewController(t)
				repository := mocks.NewMockSecretRepository(ctrl)
				repository.EXPECT().
					GetForUser(gomock.All(), gomock.All(), gomock.All()).
					Return(entity.Secret{}, nil)
				return repository
			},
			lSetup: func(t *testing.T) Logger {
				ctrl := gomock.NewController(t)
				logger := mocks.NewMockLogger(ctrl)
				logger.EXPECT().
					Error("failed to unmarhal metadata", gomock.All())
				return logger
			},
			want: want{
				err: srvErrors.ErrUnexpected,
			},
		},
		{
			name:     "repository_error",
			ctx:      goodCtx,
			secretID: 13,
			rSetup: func(t *testing.T) SecretRepository {
				ctrl := gomock.NewController(t)
				repository := mocks.NewMockSecretRepository(ctrl)
				repository.EXPECT().
					GetForUser(gomock.All(), gomock.All(), gomock.All()).
					Return(entity.Secret{}, fmt.Errorf("repository error"))
				return repository
			},
			lSetup: func(t *testing.T) Logger {
				ctrl := gomock.NewController(t)
				logger := mocks.NewMockLogger(ctrl)
				logger.EXPECT().
					Error("failed to get secret for user", gomock.All())
				return logger
			},
			want: want{
				err: srvErrors.ErrUnexpected,
			},
		},
		{
			name:     "secret_not_found",
			ctx:      goodCtx,
			secretID: 13,
			rSetup: func(t *testing.T) SecretRepository {
				ctrl := gomock.NewController(t)
				repository := mocks.NewMockSecretRepository(ctrl)
				repository.EXPECT().
					GetForUser(gomock.All(), gomock.All(), gomock.All()).
					Return(entity.Secret{}, repErrors.ErrNotFound)
				return repository
			},
			lSetup: func(t *testing.T) Logger {
				ctrl := gomock.NewController(t)
				return mocks.NewMockLogger(ctrl)
			},
			want: want{
				err: srvErrors.ErrSecretNotFound,
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			repository := test.rSetup(t)
			logger := test.lSetup(t)
			secretService := NewSecret(logger, repository)
			secret, err := secretService.Secret(test.ctx, test.secretID)
			assert.ErrorIs(t, err, test.want.err, "Retrieve secret error")
			if err == nil {
				assert.Equal(t, test.want.secret, secret, "Retrieve secret")
			}
		})
	}
}

func TestSecret_InfoList(t *testing.T) {
	userID := "1ed655b6-0738-4162-a34a-34257c0dc106"
	goodCtx := srvContext.SetUserID(context.Background(), userID)

	type want struct {
		list []dto.SecretInfo
		err  error
	}

	tests := []struct {
		name   string
		ctx    context.Context
		rSetup func(t *testing.T) SecretRepository
		lSetup func(t *testing.T) Logger
		want   want
	}{
		{
			name: "success",
			ctx:  goodCtx,
			rSetup: func(t *testing.T) SecretRepository {
				ctrl := gomock.NewController(t)
				repository := mocks.NewMockSecretRepository(ctrl)
				repository.EXPECT().
					GetAllUnencryptedByUser(gomock.All(), gomock.All()).
					Return([]entity.SecretInfo{{MetaData: "[]"}}, nil)
				return repository
			},
			lSetup: func(t *testing.T) Logger {
				ctrl := gomock.NewController(t)
				return mocks.NewMockLogger(ctrl)
			},
			want: want{
				list: []dto.SecretInfo{{Meta: []dto.MetaData{}}},
			},
		},
		{
			name: "without_user",
			ctx:  context.TODO(),
			rSetup: func(t *testing.T) SecretRepository {
				ctrl := gomock.NewController(t)
				return mocks.NewMockSecretRepository(ctrl)
			},
			lSetup: func(t *testing.T) Logger {
				ctrl := gomock.NewController(t)
				logger := mocks.NewMockLogger(ctrl)
				logger.EXPECT().
					Error("failed to get user id", gomock.All())
				return logger
			},
			want: want{
				err: srvErrors.ErrUnexpected,
			},
		},
		{
			name: "repository_error",
			ctx:  goodCtx,
			rSetup: func(t *testing.T) SecretRepository {
				ctrl := gomock.NewController(t)
				repository := mocks.NewMockSecretRepository(ctrl)
				repository.EXPECT().
					GetAllUnencryptedByUser(gomock.All(), gomock.All()).
					Return([]entity.SecretInfo{}, fmt.Errorf("repository error"))
				return repository
			},
			lSetup: func(t *testing.T) Logger {
				ctrl := gomock.NewController(t)
				logger := mocks.NewMockLogger(ctrl)
				logger.EXPECT().
					Error("failed to get secret for user", gomock.All())
				return logger
			},
			want: want{
				err: srvErrors.ErrUnexpected,
			},
		},
		{
			name: "bad_meta_data",
			ctx:  goodCtx,
			rSetup: func(t *testing.T) SecretRepository {
				ctrl := gomock.NewController(t)
				repository := mocks.NewMockSecretRepository(ctrl)
				repository.EXPECT().
					GetAllUnencryptedByUser(gomock.All(), gomock.All()).
					Return([]entity.SecretInfo{{MetaData: ""}}, nil)
				return repository
			},
			lSetup: func(t *testing.T) Logger {
				ctrl := gomock.NewController(t)
				logger := mocks.NewMockLogger(ctrl)
				logger.EXPECT().
					Error("failed to unmarhal metadata", gomock.All())
				return logger
			},
			want: want{
				err: srvErrors.ErrUnexpected,
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			repository := test.rSetup(t)
			logger := test.lSetup(t)
			secretService := NewSecret(logger, repository)

			list, err := secretService.InfoList(test.ctx)
			assert.ErrorIs(t, err, test.want.err, "Retrieve secret error")
			if err == nil {
				assert.Equal(t, test.want.list, list, "Retrieve secret info list")
			}
		})
	}
}
