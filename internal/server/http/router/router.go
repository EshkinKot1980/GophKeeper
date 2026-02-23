package router

import (
	"net/http"

	"github.com/go-chi/chi/v5"

	"github.com/EshkinKot1980/GophKeeper/internal/server/config"
	"github.com/EshkinKot1980/GophKeeper/internal/server/http/handler"
	"github.com/EshkinKot1980/GophKeeper/internal/server/http/middleware"
)

type AuthService interface {
	handler.AuthService
	middleware.AuthService
}

type Logger interface {
	handler.Logger
	middleware.HTTPloger
}

type SecretService = handler.SecretService

// NewRouter инициализирует хендлеры и создает роутер *chiMux
func NewRouter(cfg *config.Config, l Logger, a AuthService, s SecretService) http.Handler {
	authorizer := middleware.NewAuthorizer(a)
	logger := middleware.NewLogger(l)
	authHandler := handler.NewAuth(a, l, cfg.AuthBodyMaxSize)
	secretHandler := handler.NewSecret(s, l)

	router := chi.NewRouter()

	router.Route("/api", func(r chi.Router) {
		r.Use(logger.Log)
		r.Route("/register", func(r chi.Router) {
			r.Post("/", authHandler.Register)
		})
		r.Route("/login", func(r chi.Router) {
			r.Post("/", authHandler.Login)
		})

		r.Group(func(r chi.Router) {
			r.Use(authorizer.Authorize)

			r.Route("/secret", func(r chi.Router) {
				r.Post("/", secretHandler.Upload)
				r.Get("/{id}", secretHandler.Get)
				r.Get("/", secretHandler.List)
			})
		})
	})

	return router
}
