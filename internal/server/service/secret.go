// Пакет service содержит сервисный слой серверной части приложения
package service

import (
	"context"
	"encoding/json"
	"errors"

	"github.com/EshkinKot1980/GophKeeper/internal/common/dto"
	"github.com/EshkinKot1980/GophKeeper/internal/server/entity"
	"github.com/EshkinKot1980/GophKeeper/internal/server/http/middleware"
	repErrors "github.com/EshkinKot1980/GophKeeper/internal/server/repository/errors"
	srvErrors "github.com/EshkinKot1980/GophKeeper/internal/server/service/errors"
)

type SecretRepository interface {
	// Create создает пользовательский секрет в БД.
	Create(ctx context.Context, secret entity.Secret) error
	// GetForUser возвращает пользовательский секрет по secretID и userID.
	GetForUser(ctx context.Context, secretID uint64, userID string) (entity.Secret, error)
}

// Secret сервис загрузки и отдачи секретов пользователя
type Secret struct {
	logger     Logger
	repository SecretRepository
}

func NewSecret(l Logger, s SecretRepository) *Secret {
	return &Secret{logger: l, repository: s}
}

// Save сохраняет секрет на сервере
func (s *Secret) Save(ctx context.Context, secret *dto.SecretRequest) error {
	userID, ok := ctx.Value(middleware.KeyUserID).(string)
	if !ok {
		s.logger.Error("failed to get user id", srvErrors.ErrUnexpected)
		return srvErrors.ErrUnexpected
	}

	enity := entity.Secret{
		UserID:        userID,
		DataType:      secret.DataType,
		Name:          secret.Name,
		MetaData:      "[]", //TODO: ...
		EncryptedKey:  secret.EncrData.Key,
		EncryptedData: secret.EncrData.Data,
	}

	err := s.repository.Create(ctx, enity)
	if err != nil {
		s.logger.Error("failed create secret", err)
		return srvErrors.ErrUnexpected
	}

	return nil
}

// Secret возвращает секрет по secretID, если он принадлежит текущему пользователю.
func (s *Secret) Secret(ctx context.Context, secretID uint64) (dto.SecretResponse, error) {
	userID, ok := ctx.Value(middleware.KeyUserID).(string)
	if !ok {
		s.logger.Error("failed to get user id", srvErrors.ErrUnexpected)
		return dto.SecretResponse{}, srvErrors.ErrUnexpected
	}

	entity, err := s.repository.GetForUser(ctx, secretID, userID)
	if err != nil {
		if errors.Is(err, repErrors.ErrNotFound) {
			return dto.SecretResponse{}, srvErrors.ErrSecretNotFound
		}
		s.logger.Error("failed to get secret for user", err)
		return dto.SecretResponse{}, srvErrors.ErrUnexpected
	}

	var meta []dto.MetaData
	err = json.Unmarshal([]byte(entity.MetaData), &meta)
	if err != nil {
		s.logger.Error("failed to unmarhal metadata", err)
		return dto.SecretResponse{}, srvErrors.ErrUnexpected
	}

	return dto.SecretResponse{
		ID:       entity.ID,
		DataType: entity.DataType,
		Name:     entity.Name,
		Meta:     meta,
		Created:  entity.Created,
		Updated:  entity.Updated,
		EncrData: dto.EncryptedData{
			Key:  entity.EncryptedKey,
			Data: entity.EncryptedData,
		},
	}, nil
}
