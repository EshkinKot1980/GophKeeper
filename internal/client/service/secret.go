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
		return fmt.Errorf("unauthorized, try login to system :%w", err)
	}

	secret.EncrData, err = ecryptData(masterKey, data)
	if err != nil {
		return fmt.Errorf("failed to ecrypt secret :%w", err)
	}

	token, err := s.storage.Token()
	if err != nil {
		return fmt.Errorf("unauthorized, try login to system :%w", err)
	}

	return s.client.Upload(secret, token)
}

func ecryptData(masterKey, payload []byte) (*dto.EncryptedData, error) {
	var result dto.EncryptedData

	key, err := crypto.GenerateRandomBytes(32)
	if err != nil {
		return nil, fmt.Errorf("failed to generate DEK: %w", err)
	}

	result.Data, err = crypto.EncryptAES(key, payload)
	if err != nil {
		return nil, fmt.Errorf("failed to encrypt payload: %w", err)
	}

	encryptedKey, err := crypto.EncryptAES(masterKey, key)
	if err != nil {
		return nil, fmt.Errorf("failed to encrypt DEK: %w", err)
	}
	result.Key = base64.RawStdEncoding.EncodeToString(encryptedKey)

	return &result, nil
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
		return nil, fmt.Errorf("failed to decrypt data content: %w", err)
	}

	return decryptedData, nil
}
