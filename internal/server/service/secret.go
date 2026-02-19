// Пакет service содержит сервисный слой серверной части приложения
package service

import (
	"context"

	"github.com/EshkinKot1980/GophKeeper/internal/common/dto"
	"github.com/EshkinKot1980/GophKeeper/internal/server/entity"
	"github.com/EshkinKot1980/GophKeeper/internal/server/http/middleware"
	srvErrors "github.com/EshkinKot1980/GophKeeper/internal/server/service/errors"
)

type SecretRepository interface {
	// Create создает запись секрета
	Create(ctx context.Context, secret entity.Secret) error
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
