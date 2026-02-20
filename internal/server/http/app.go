package http

import (
	"context"
	"log"
	"net/http"
	"time"

	"github.com/EshkinKot1980/GophKeeper/internal/server/config"
	"github.com/EshkinKot1980/GophKeeper/internal/server/http/handler"
	"github.com/EshkinKot1980/GophKeeper/internal/server/http/middleware"
	"github.com/EshkinKot1980/GophKeeper/internal/server/logger"
	"github.com/go-chi/chi/v5"
)

type AuthService interface {
	handler.AuthService
	middleware.AuthService
}

type SecretService = handler.SecretService

type App struct {
	config        *config.Config
	logger        *logger.Logger
	authService   AuthService
	secretService SecretService
}

func NewApp(c *config.Config, l *logger.Logger, a AuthService, s SecretService) *App {
	return &App{config: c, logger: l, authService: a, secretService: s}
}

func (a *App) Run(ctx context.Context) error {
	srv := &http.Server{Addr: a.config.HTTPSaddr, Handler: a.newRouter()}
	errChan := make(chan error)

	go func() {
		err := srv.ListenAndServeTLS(a.config.TLScert, a.config.TLSkey)
		if err != nil && err != http.ErrServerClosed {
			errChan <- err
		}
	}()

	select {
	case err := <-errChan:
		return err
	case <-time.After(time.Second):
		log.Printf("server listening on %s\n", a.config.HTTPSaddr)
	}

	<-ctx.Done()
	log.Println("shutting down http server gracefully")
	shutdownCtx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer func() {
		log.Println("http server stopped")
		cancel()
	}()

	return srv.Shutdown(shutdownCtx)
}

func (a *App) newRouter() http.Handler {
	authorizer := middleware.NewAuthorizer(a.authService)
	logger := middleware.NewLogger(a.logger)
	authHandler := handler.NewAuth(a.authService, a.logger)
	secretHandler := handler.NewSecret(a.secretService, a.logger)

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
			})
		})
	})

	return router
}
