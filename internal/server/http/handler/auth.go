package handler

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"

	"github.com/EshkinKot1980/GophKeeper/internal/common/dto"
	srvErrors "github.com/EshkinKot1980/GophKeeper/internal/server/service/errors"
)

type AuthService interface {
	// Register регистрирует пользователя по логину и паролю.
	Register(ctx context.Context, c dto.Credentials) (dto.AuthResponse, error)
	// Login выполняет вход пользователя в систему по логину с паролем.
	Login(ctx context.Context, c dto.Credentials) (dto.AuthResponse, error)
}

// Auth обработчик запросов регистрации и логина
type Auth struct {
	service     AuthService
	logger      Logger
	bodyMaxSize int64
}

func NewAuth(srv AuthService, l Logger, bodyMaxSize int64) *Auth {
	return &Auth{service: srv, logger: l, bodyMaxSize: bodyMaxSize}
}

// Register регистрирует пользователя по логину и паролю.
// В случае успеха, возвращает JSON, содержащий токен (JWT)
// и соль для создания мастер ключа, закодированную base64
func (h *Auth) Register(w http.ResponseWriter, r *http.Request) {
	var credentials dto.Credentials

	body := http.MaxBytesReader(w, r.Body, h.bodyMaxSize)
	if err := json.NewDecoder(body).Decode(&credentials); err != nil {
		http.Error(w, "invalid credentials format", http.StatusBadRequest)
		return
	}

	resp, err := h.service.Register(r.Context(), credentials)
	if err != nil {
		switch {
		case errors.Is(err, srvErrors.ErrAuthInvalidCredentials):
			http.Error(w, err.Error(), http.StatusBadRequest)
		case errors.Is(err, srvErrors.ErrAuthUserAlreadyExists):
			http.Error(w, err.Error(), http.StatusConflict)
		default:
			http.Error(w, statusText500, http.StatusInternalServerError)
		}

		return
	}

	newJSONwriter(w, h.logger).write(resp, "register response", http.StatusOK)
}

// Login вход пользователя в систему по логину с паролем.
// В случае успеха, возвращает JSON, содержащий токен (JWT)
// и соль для создания мастер ключа, закодированную base64
func (h *Auth) Login(w http.ResponseWriter, r *http.Request) {
	var credentials dto.Credentials

	body := http.MaxBytesReader(w, r.Body, h.bodyMaxSize)
	if err := json.NewDecoder(body).Decode(&credentials); err != nil {
		http.Error(w, "invalid credentials format", http.StatusBadRequest)
		return
	}

	resp, err := h.service.Login(r.Context(), credentials)
	if err != nil {
		if errors.Is(err, srvErrors.ErrAuthInvalidCredentials) {
			http.Error(w, "", http.StatusUnauthorized)
		} else {
			http.Error(w, statusText500, http.StatusInternalServerError)
		}
		return
	}

	newJSONwriter(w, h.logger).write(resp, "register response", http.StatusOK)
}
