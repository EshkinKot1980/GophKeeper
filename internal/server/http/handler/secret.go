// Пакет handler содержит обработчики http запросов
package handler

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"github.com/EshkinKot1980/GophKeeper/internal/common/dto"
	srvErrors "github.com/EshkinKot1980/GophKeeper/internal/server/service/errors"
)

type SecretService interface {
	// Save сохраняет секрет на сервере
	Save(ctx context.Context, secret *dto.SecretRequest) error
	// Secret возвращает секрет по secretID, если он принадлежит текущему пользователю.
	Secret(ctx context.Context, secretID uint64) (dto.SecretResponse, error)
	// InfoList возвращает информаци о всех секретах пользователя.
	InfoList(ctx context.Context) ([]dto.SecretInfo, error)
}

// Secret обработчик запросов загрузки и отдачи секретов пользователя
type Secret struct {
	service SecretService
	logger  Logger
}

func NewSecret(srv SecretService, l Logger) *Secret {
	return &Secret{service: srv, logger: l}
}

// Upload загружает данные секрета
func (s *Secret) Upload(w http.ResponseWriter, r *http.Request) {
	var secret dto.SecretRequest

	if err := json.NewDecoder(r.Body).Decode(&secret); err != nil {
		http.Error(w, "invalid request format", http.StatusBadRequest)
		return
	}

	err := s.service.Save(r.Context(), &secret)
	if err != nil {
		if errors.Is(err, srvErrors.ErrSecretInvalidData) {
			http.Error(w, err.Error(), http.StatusBadRequest)
		} else {
			http.Error(w, statusText500, http.StatusInternalServerError)
		}
		return
	}
}

// Get отдает секрет по id, который берет из пути.
func (s *Secret) Get(w http.ResponseWriter, r *http.Request) {
	secretID, err := strconv.ParseUint(r.PathValue("id"), 10, 64)
	if err != nil {
		http.Error(w, "invalid secret id", http.StatusBadRequest)
		return
	}

	secret, err := s.service.Secret(r.Context(), secretID)

	if err != nil {
		if errors.Is(err, srvErrors.ErrSecretNotFound) {
			http.Error(w, "", http.StatusNotFound)
		} else {
			http.Error(w, statusText500, http.StatusInternalServerError)
		}
		return
	}

	newJSONwriter(w, s.logger).write(secret, "secret", http.StatusOK)
}

// InfoList возвращает информацию о всех секретах пользователя.
func (s *Secret) List(w http.ResponseWriter, r *http.Request) {
	list, err := s.service.InfoList(r.Context())
	if err != nil {
		http.Error(w, statusText500, http.StatusInternalServerError)
		return
	}

	newJSONwriter(w, s.logger).write(list, "secret info list", http.StatusOK)
}
