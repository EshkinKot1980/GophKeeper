// Пакет storage предоставляет локальное хранилище данных на клиенте.
package storage

import (
	"fmt"
	"os"
	"path/filepath"
)

const (
	parentDirName  = ".gophkeeper"
	storageDirName = "cache"
	tokenFileName  = "token"
	keyFileName    = "key"
)

// Файловое хранилище данных для токена авторизации и мастер ключа.
// Не самое безопасное решение, в дальнейшем планирую переделать на go-keyring
// или на хранилище в памяти, в случае если пределаю клиент на полностью интерактивный cli.
type FileStorage struct {
	tokenPath string
	keyPath   string
}

func NewFileSorage() (*FileStorage, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return nil, fmt.Errorf("failed to get user home dir: %w", err)
	}

	path := filepath.Join(homeDir, parentDirName, storageDirName)
	// Создаем директорию хранилища, если её нет
	if err := os.MkdirAll(path, 0700); err != nil {
		return nil, fmt.Errorf("failed to create storage dir: %w", err)
	}

	s := FileStorage{
		tokenPath: filepath.Join(path, tokenFileName),
		keyPath:   filepath.Join(path, keyFileName),
	}

	return &s, nil
}

// PutToken сохранение токена
func (s *FileStorage) PutToken(token string) error {
	err := os.WriteFile(s.tokenPath, []byte(token), 0600)
	if err != nil {
		return fmt.Errorf("failed to write token: %w", err)
	}
	return nil
}

// Token получение токена
func (s *FileStorage) Token() (string, error) {
	token, err := os.ReadFile(s.tokenPath)
	if err != nil {
		return "", fmt.Errorf("failed to write token: %w", err)
	}
	return string(token), nil
}

// PutKey сохранение ключа
func (s *FileStorage) PutKey(key []byte) error {
	err := os.WriteFile(s.keyPath, key, 0600)
	if err != nil {
		return fmt.Errorf("failed to write key: %w", err)
	}
	return nil
}

// Key получение ключа
func (s *FileStorage) Key() ([]byte, error) {
	key, err := os.ReadFile(s.tokenPath)
	if err != nil {
		return []byte{}, fmt.Errorf("failed to write token: %w", err)
	}
	return key, nil
}
