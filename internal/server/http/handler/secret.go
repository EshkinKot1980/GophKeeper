// Пакет handler содержит обработчики http запросов
package handler

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"

	"github.com/EshkinKot1980/GophKeeper/internal/common/dto"
	srvErrors "github.com/EshkinKot1980/GophKeeper/internal/server/service/errors"
)

type SecretService interface {
	// Save сохраняет секрет на сервере
	Save(ctx context.Context, secret *dto.SecretRequest) error
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
		if errors.Is(err, srvErrors.ErrSecretInalidData) {
			http.Error(w, err.Error(), http.StatusBadRequest)
		} else {
			http.Error(w, statusText500, http.StatusInternalServerError)
		}
		return
	}
}
