package main

import (
	"context"
	"fmt"
	"log"
	"os/signal"
	"syscall"

	"github.com/EshkinKot1980/GophKeeper/internal/server/config"
	"github.com/EshkinKot1980/GophKeeper/internal/server/http"
	"github.com/EshkinKot1980/GophKeeper/internal/server/logger"
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

	fmt.Println(cfg)

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	defer stop()

	// db, err := pg.NewDB(ctx, cfg.DatabaseDSN)
	// if err != nil {
	// 	return fmt.Errorf("failed to init DB: %w", err)
	// }
	// defer db.Close()

	logger, err := logger.New()
	if err != nil {
		return fmt.Errorf("failed to init logger: %w", err)
	}
	defer logger.Sync()

	httpsServer := http.NewApp(cfg, logger)
	return httpsServer.Run(ctx)
}
