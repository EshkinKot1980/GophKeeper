// Пакет service содержит сервисный слой клиентской части приложения
package service

import "github.com/EshkinKot1980/GophKeeper/internal/common/dto"

// Storage хранилище данных для токена авторизации и мастер ключа.
type Storage interface {
	// PutToken сохранение токена
	PutToken(token string) error
	// Token получение токена
	Token() (string, error)
	// PutKey сохранение ключа
	PutKey(key []byte) error
	// Key получение ключа
	Key() ([]byte, error)
}

// Client клиент для взаимодействия с сервером
type Client interface {
	// Register регистрирует пользователя в системе
	Register(cr dto.Credentials) (dto.AuthResponse, error)
	// Login осуществляет вход пользователя в систему
	Login(cr dto.Credentials) (dto.AuthResponse, error)
	// Upload coхраняет секрет на сервере
	Upload(data dto.SecretRequest, token string) error
	// Retrieve получает секрет пользователя с ервера
	Retrieve(id uint64, token string) (dto.SecretResponse, error)
	// InfoList получает информацию о всех секретах пользователя с сервера.
	InfoList(token string) ([]dto.SecretInfo, error)
}
