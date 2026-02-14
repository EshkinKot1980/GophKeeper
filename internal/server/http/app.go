package http

import (
	"context"
	"log"
	"net/http"
	"time"

	"github.com/EshkinKot1980/GophKeeper/internal/server/config"
	"github.com/EshkinKot1980/GophKeeper/internal/server/logger"
	"github.com/go-chi/chi/v5"
)

type App struct {
	config *config.Config
	logger *logger.Logger
}

func NewApp(c *config.Config, l *logger.Logger) *App {
	return &App{config: c, logger: l}
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
	r := chi.NewRouter()
	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("welcome"))
	})
	return r
}
