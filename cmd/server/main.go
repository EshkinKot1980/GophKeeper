package main

import (
	"context"
	"fmt"
	"log"
	"os/signal"
	"syscall"

	"github.com/EshkinKot1980/GophKeeper/internal/common/crypto"
	"github.com/EshkinKot1980/GophKeeper/internal/server/config"
	"github.com/EshkinKot1980/GophKeeper/internal/server/http"
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
	authService := service.NewAuth(userRepository, logger, jwtPublicKey, jwtPrivateKey)

	secretRepository := repository.NewSecret(db)
	secretService := service.NewSecret(logger, secretRepository)

	httpsServer := http.NewApp(cfg, logger, authService, secretService)
	return httpsServer.Run(ctx)
}
