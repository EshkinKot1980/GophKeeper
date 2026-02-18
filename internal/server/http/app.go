package http

import (
	"context"
	"log"
	"net/http"
	"time"

	"github.com/EshkinKot1980/GophKeeper/internal/common/dto"
	"github.com/EshkinKot1980/GophKeeper/internal/server/config"
	"github.com/EshkinKot1980/GophKeeper/internal/server/entity"
	"github.com/EshkinKot1980/GophKeeper/internal/server/http/handler"
	"github.com/EshkinKot1980/GophKeeper/internal/server/http/middleware"
	"github.com/EshkinKot1980/GophKeeper/internal/server/logger"
	"github.com/go-chi/chi/v5"
)

type AuthService interface {
	Register(ctx context.Context, c dto.Credentials) (dto.AuthResponse, error)
	Login(ctx context.Context, c dto.Credentials) (dto.AuthResponse, error)
	User(ctx context.Context, token string) (entity.User, error)
}

type App struct {
	config *config.Config
	logger *logger.Logger
	auth   AuthService
}

func NewApp(c *config.Config, l *logger.Logger, a AuthService) *App {
	return &App{config: c, logger: l, auth: a}
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
	authorizer := middleware.NewAuthorizer(a.auth)
	authHandler := handler.NewAuth(a.auth, a.logger)

	router := chi.NewRouter()

	router.Route("/api", func(r chi.Router) {
		r.Route("/register", func(r chi.Router) {
			r.Post("/", authHandler.Register)
		})
		r.Route("/login", func(r chi.Router) {
			r.Post("/", authHandler.Login)
		})

		r.Group(func(r chi.Router) {
			r.Use(authorizer.Authorize)

			r.Get("/", func(w http.ResponseWriter, r *http.Request) {
				w.Write([]byte("welcome"))
			})
		})
	})

	router.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("welcome"))
	})

	return router
}
