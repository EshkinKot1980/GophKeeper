// Пакет service содержит сервисный слой клиентской части приложения
package service

import (
	"encoding/base64"
	"fmt"

	"github.com/EshkinKot1980/GophKeeper/internal/common/crypto"
	"github.com/EshkinKot1980/GophKeeper/internal/common/dto"
)

// Auth сервис для регистрации аутентификации и авторизации
type Auth struct {
	client  Client
	storage Storage
}

func NewAuth(c Client, s Storage) *Auth {
	return &Auth{client: c, storage: s}
}

// Register регистрирует пользователя в системе
func (a *Auth) Register(cr dto.Credentials) error {
	resp, err := a.client.Register(cr)
	if err != nil {
		return err
	}

	return a.storeAuthData(cr, resp)
}

// Login осуществляет вход пользователя в систему
func (a *Auth) Login(cr dto.Credentials) error {
	resp, err := a.client.Login(cr)
	if err != nil {
		return err
	}

	return a.storeAuthData(cr, resp)
}

// storeAuthData вычисляет мастел ключ
// и сохвраняет его вместе с токеном в локальное хранилище.
func (a *Auth) storeAuthData(cr dto.Credentials, resp dto.AuthResponse) error {
	salt, err := base64.RawStdEncoding.DecodeString(resp.EncrSalt)
	if err != nil {
		return fmt.Errorf("the server returned invalid data")
	}

	// Вычисляем  мастер ключ для шифрования данных
	masterKey, err := crypto.DeriveKey([]byte(cr.Password), salt)
	if err != nil {
		return fmt.Errorf("failed to derive key")
	}

	// сохраняем ключ в хранилище
	err = a.storage.PutKey(masterKey)
	if err != nil {
		return fmt.Errorf("failed to store key")
	}

	// сохраняем токен в хранилище
	err = a.storage.PutToken(resp.Token)
	if err != nil {
		return fmt.Errorf("failed to store token")
	}

	return nil
}
