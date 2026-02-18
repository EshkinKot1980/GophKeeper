// Пакет cli содержит интерфейс омандной строки.
package cli

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/EshkinKot1980/GophKeeper/internal/client/config"
	"github.com/EshkinKot1980/GophKeeper/internal/client/http"
	"github.com/EshkinKot1980/GophKeeper/internal/client/service"
	"github.com/EshkinKot1980/GophKeeper/internal/client/storage"
	"github.com/EshkinKot1980/GophKeeper/internal/common/dto"
)

// Auth сервис для регистрации аутентификации и авторизации
type AuthService interface {
	//Register регистрирует пользователя в системе
	Register(cr dto.Credentials) error
	//Login осуществляет вход пользователя в систему
	Login(cr dto.Credentials) error
}

var (
	cfg         *config.Config
	httpClient  *http.Client
	fileStorage *storage.FileStorage
	authService AuthService
)

// Корневая команда приложения Cobra
var rootCmd = &cobra.Command{
	Use:   "gophkeeper",
	Short: "GophKeeper CLI client",
	Long:  `Secure client for store passwords, credit cards data, and binary files.`,

	// PersistentPreRunE выполняется перед любой командой
	// Здесь мы инициализируем подключение к БД и загружаем конфиг.
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		var err error
		cfg, err = config.Load()
		if err != nil {
			return fmt.Errorf("failed to load config: %w", err)
		}

		httpClient = http.NewClient(cfg.ServerAddr, cfg.AllowSelfSignedCert)

		fileStorage, err = storage.NewFileSorage()
		if err != nil {
			return fmt.Errorf("failed init storage: %w", err)
		}

		authService = service.NewAuth(httpClient, fileStorage)

		return nil
	},
}

// Execute - точка входа для CLI
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		// Cobra сама выводит ошибку
		os.Exit(1)
	}
}
