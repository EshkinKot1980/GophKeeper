package dto

import (
	"fmt"
	"unicode"
)

const (
	CredentialsLoginMaxLen    = 64
	CredentialsLoginMinLen    = 3
	CredentialsPasswordMinLen = 8
)

type AuthResponse struct {
	// токен для авторизации (JWT)
	Token string `json:"token"`
	// Соль для создания мастер из пароля ключа закодированная base64
	EncrSalt string `json:"encr_salt"`
}

type Credentials struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}

// Validate проверяет учетные данные пользователя при регистрации,
// используется на сервере
func (cr Credentials) Validate() error {
	if err := cr.ValidateLogin(); err != nil {
		return err
	}

	if err := cr.ValidatePassword(); err != nil {
		return err
	}

	return nil
}

// ValidateLogin проверяет логин пользователя при регистрации, используется на клиенте
func (cr Credentials) ValidateLogin() error {
	if len(cr.Login) < CredentialsLoginMinLen {
		return fmt.Errorf("login too short (min %d chars)", CredentialsLoginMinLen)
	}

	if len(cr.Login) > CredentialsLoginMaxLen {
		return fmt.Errorf("login too long (max %d chars)", CredentialsLoginMaxLen)
	}

	var hasLetter, hasDigit bool
	for _, c := range cr.Login {
		if unicode.IsLetter(c) {
			hasLetter = true
		}
		if unicode.IsDigit(c) {
			hasDigit = true
		}
	}

	if !(hasLetter || hasDigit) {
		return fmt.Errorf("login must contain letters or digits")
	}

	return nil
}

// ValidatePassword проверяет пароль пользователя при регистрации, используется на клиенте
func (cr Credentials) ValidatePassword() error {
	if len(cr.Password) < CredentialsPasswordMinLen {
		return fmt.Errorf("password too short (min %d chars)", CredentialsPasswordMinLen)
	}

	var hasLetter, hasDigit bool
	for _, c := range cr.Password {
		if unicode.IsLetter(c) {
			hasLetter = true
		}
		if unicode.IsDigit(c) {
			hasDigit = true
		}
	}

	if !hasLetter || !hasDigit {
		return fmt.Errorf("password must contain both letters and digits")
	}

	return nil
}
