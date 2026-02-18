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
	//Register регистрирует пользователя в системе
	Register(cr dto.Credentials) (dto.AuthResponse, error)
	//Login осуществляет вход пользователя в систему
	Login(cr dto.Credentials) (dto.AuthResponse, error)
}
