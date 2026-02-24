package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os/signal"
	"syscall"
	"time"

	"github.com/EshkinKot1980/GophKeeper/internal/common/crypto"
	"github.com/EshkinKot1980/GophKeeper/internal/server/config"
	"github.com/EshkinKot1980/GophKeeper/internal/server/http/router"
	"github.com/EshkinKot1980/GophKeeper/internal/server/logger"
	"github.com/EshkinKot1980/GophKeeper/internal/server/repository"
	"github.com/EshkinKot1980/GophKeeper/internal/server/repository/pg"
	"github.com/EshkinKot1980/GophKeeper/internal/server/service"
)

func main() {
	if err := run(); err != nil {
		log.Fatal(err)
	}
}

func run() error {
	cfg, err := config.Load()
	if err != nil {
		return fmt.Errorf("failed to load config: %w", err)
	}

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	defer stop()

	db, err := pg.NewDB(ctx, cfg.DatabaseDSN)
	if err != nil {
		return fmt.Errorf("failed to init DB: %w", err)
	}
	defer db.Close()

	logger, err := logger.New()
	if err != nil {
		return fmt.Errorf("failed to init logger: %w", err)
	}
	defer logger.Sync()

	jwtPrivateKey, err := crypto.LoadPrivateKey(cfg.JWTpriv)
	if err != nil {
		return fmt.Errorf("failed to load jwt private key: %w", err)
	}
	jwtPublicKey, err := crypto.LoadPublicKey(cfg.JWTpub)
	if err != nil {
		return fmt.Errorf("failed to load jwt ppublic key: %w", err)
	}

	userRepository := repository.NewUser(db)
	authService := service.NewAuth(userRepository, logger, jwtPublicKey, jwtPrivateKey, cfg.TokenTTL)

	secretRepository := repository.NewSecret(db)
	secretService := service.NewSecret(logger, secretRepository)

	router := router.NewRouter(cfg, logger, authService, secretService)
	return sevreHTTPS(ctx, cfg, logger, router)
}

func sevreHTTPS(ctx context.Context, cfg *config.Config, logger *logger.Logger, router http.Handler) error {
	srv := &http.Server{Addr: cfg.HTTPSaddr, Handler: router}
	errChan := make(chan error, 1)

	go func() {
		err := srv.ListenAndServeTLS(cfg.TLScert, cfg.TLSkey)
		if err != nil && err != http.ErrServerClosed {
			errChan <- err
		}
	}()

	select {
	case err := <-errChan:
		return err
	case <-time.After(time.Second):
		logger.Info("server listening on " + cfg.HTTPSaddr)
	}

	<-ctx.Done()
	logger.Info("shutting down http server gracefully")
	shutdownCtx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer func() {
		logger.Info("http server stopped")
		cancel()
	}()

	return srv.Shutdown(shutdownCtx)
}
