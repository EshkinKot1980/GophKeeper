// Пакет service содержит сервисный слой серверной части приложения
package service

import (
	"context"
	"encoding/json"
	"errors"

	"github.com/EshkinKot1980/GophKeeper/internal/common/dto"
	"github.com/EshkinKot1980/GophKeeper/internal/server/entity"
	repErrors "github.com/EshkinKot1980/GophKeeper/internal/server/repository/errors"
	srvContext "github.com/EshkinKot1980/GophKeeper/internal/server/service/context"
	srvErrors "github.com/EshkinKot1980/GophKeeper/internal/server/service/errors"
)

type SecretRepository interface {
	// Create создает пользовательский секрет в БД.
	Create(ctx context.Context, secret entity.Secret) error
	// GetForUser возвращает пользовательский секрет по secretID и userID.
	GetForUser(ctx context.Context, secretID uint64, userID string) (entity.Secret, error)
	// GetAlluUnencryptedByUser возвращает не зашифрованные данные для всех записей пользователя
	GetAllUnencryptedByUser(ctx context.Context, userID string) ([]entity.SecretInfo, error)
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
	userID, err := srvContext.UserID(ctx)
	if err != nil {
		s.logger.Error("failed to get user id", err)
		return srvErrors.ErrUnexpected
	}

	meta, err := json.Marshal(secret.Meta)
	if err != nil {
		s.logger.Error("failed encode secret metadata to json", err)
		return srvErrors.ErrUnexpected
	}

	enity := entity.Secret{
		UserID:        userID,
		DataType:      secret.DataType,
		Name:          secret.Name,
		MetaData:      string(meta),
		EncryptedKey:  secret.EncrData.Key,
		EncryptedData: secret.EncrData.Data,
	}

	err = s.repository.Create(ctx, enity)
	if err != nil {
		s.logger.Error("failed create secret", err)
		return srvErrors.ErrUnexpected
	}

	return nil
}

// Secret возвращает секрет по secretID, если он принадлежит текущему пользователю.
func (s *Secret) Secret(ctx context.Context, secretID uint64) (dto.SecretResponse, error) {
	userID, err := srvContext.UserID(ctx)
	if err != nil {
		s.logger.Error("failed to get user id", err)
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

// InfoList возвращает информацию о всех секретах пользователя.
func (s *Secret) InfoList(ctx context.Context) ([]dto.SecretInfo, error) {
	userID, err := srvContext.UserID(ctx)
	if err != nil {
		s.logger.Error("failed to get user id", err)
		return nil, srvErrors.ErrUnexpected
	}

	secrets, err := s.repository.GetAllUnencryptedByUser(ctx, userID)
	if err != nil {
		s.logger.Error("failed to get secret for user", err)
		return nil, srvErrors.ErrUnexpected
	}

	var list []dto.SecretInfo
	for _, secret := range secrets {
		var meta []dto.MetaData
		err = json.Unmarshal([]byte(secret.MetaData), &meta)
		if err != nil {
			s.logger.Error("failed to unmarhal metadata", err)
			return nil, srvErrors.ErrUnexpected
		}

		list = append(
			list,
			dto.SecretInfo{
				ID:       secret.ID,
				DataType: secret.DataType,
				Name:     secret.Name,
				Meta:     meta,
				Created:  secret.Created,
			},
		)
	}
	return list, nil
}
