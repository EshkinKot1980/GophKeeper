// Пакет service содержит сервисный слой клиентской части приложения
package service

import (
	"encoding/base64"
	"fmt"

	"github.com/EshkinKot1980/GophKeeper/internal/common/crypto"
	"github.com/EshkinKot1980/GophKeeper/internal/common/dto"
)

// Secret сервис для работы с секретными данными пользователя
type Secret struct {
	client  Client
	storage Storage
}

func NewSecret(c Client, s Storage) *Secret {
	return &Secret{client: c, storage: s}
}

// Upload отправляет данные на сервер.
// Принимает частино заполненный dto.SecretRequest и данные,
// которые нужно зашифровать.
func (s *Secret) Upload(secret dto.SecretRequest, data []byte) error {
	masterKey, err := s.storage.Key()
	if err != nil {
		return fmt.Errorf("authorization failed :%w", err)
	}

	secret.EncrData, err = ecryptData(masterKey, data)
	if err != nil {
		return fmt.Errorf("failed to ecrypt secret :%w", err)
	}

	token, err := s.storage.Token()
	if err != nil {
		return fmt.Errorf("authorization failed :%w", err)
	}

	return s.client.Upload(secret, token)
}

// GetSecretAndInfo получает секрет пользователя с сервера по id,
// возвращает расшиврованные данные в виде []byte и  информацию о секрете
func (s *Secret) GetSecretAndInfo(id uint64) ([]byte, dto.SecretInfo, error) {
	var info dto.SecretInfo

	masterKey, err := s.storage.Key()
	if err != nil {
		return nil, info, fmt.Errorf("authorization failed :%w", err)
	}
	token, err := s.storage.Token()
	if err != nil {
		return nil, info, fmt.Errorf("authorization failed :%w", err)
	}

	resp, err := s.client.Retrieve(id, token)
	if err != nil {
		return nil, info, err
	}

	secret, err := deryptData(masterKey, &resp.EncrData)
	if err != nil {
		return nil, info, fmt.Errorf("failed to decrypt secret :%w", err)
	}

	info = dto.SecretInfo{
		ID:       resp.ID,
		DataType: resp.DataType,
		Name:     resp.Name,
		Meta:     resp.Meta,
		Created:  resp.Created,
	}

	return secret, info, nil
}

// InfoList получает информацию о всех секретах пользователя с сервера.
func (s *Secret) InfoList() ([]dto.SecretInfo, error) {
	token, err := s.storage.Token()
	if err != nil {
		return nil, fmt.Errorf("authorization failed :%w", err)
	}

	list, err := s.client.InfoList(token)
	if err != nil {
		return nil, err
	}
	return list, nil
}

func ecryptData(masterKey, payload []byte) (dto.EncryptedData, error) {
	var result dto.EncryptedData

	key, err := crypto.GenerateRandomBytes(32)
	if err != nil {
		return result, fmt.Errorf("failed to generate DEK: %w", err)
	}

	result.Data, err = crypto.EncryptAES(key, payload)
	if err != nil {
		return result, fmt.Errorf("failed to encrypt payload: %w", err)
	}

	encryptedKey, err := crypto.EncryptAES(masterKey, key)
	if err != nil {
		return result, fmt.Errorf("failed to encrypt DEK: %w", err)
	}
	result.Key = base64.RawStdEncoding.EncodeToString(encryptedKey)

	return result, nil
}

func deryptData(masterKey []byte, data *dto.EncryptedData) ([]byte, error) {
	if data == nil {
		return nil, fmt.Errorf("the server returned invalid data: EncryptedData is nil")
	}

	encryptedKey, err := base64.RawStdEncoding.DecodeString(data.Key)
	if err != nil {
		return nil, fmt.Errorf("the server returned invalid data: bad key")
	}

	key, err := crypto.DecryptAES(masterKey, encryptedKey)
	if err != nil {
		return nil, fmt.Errorf("failed to decrypt DEK (invalid master key?): %w", err)
	}

	decryptedData, err := crypto.DecryptAES(key, data.Data)
	if err != nil {
		return nil, fmt.Errorf("failed to decrypt data: %w", err)
	}

	return decryptedData, nil
}
