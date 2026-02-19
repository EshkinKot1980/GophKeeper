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

// SecretService сервис для работы с секретными данными пользователя
type SecretService interface {
	// Upload отправляет зашифрованные данные на сервер.
	// Принимает частино заполненный dto.SecretRequest и данные, которые нужно зашифровать.
	Upload(secret dto.SecretRequest, data *dto.PlainData) error
}

var (
	cfg           *config.Config
	authService   AuthService
	secretServise SecretService
)

// Корневая команда приложения Cobra
var rootCmd = &cobra.Command{
	Use:   "gophkeeper",
	Short: "GophKeeper CLI client",
	Long:  `Secure client for store passwords, credit cards data, and binary files.`,

	// PersistentPreRunE выполняется перед любой командой
	// Здесь мы загружаем конфиг и инициализируем сервисы.
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		var err error
		cfg, err = config.Load()
		if err != nil {
			return fmt.Errorf("failed to load config: %w", err)
		}

		baseURL := http.Scheme + cfg.ServerAddr + http.APIprefix
		httpClient := http.NewClient(baseURL, cfg.AllowSelfSignedCert)

		fileStorage, err := storage.NewFileSorage()
		if err != nil {
			return fmt.Errorf("failed init storage: %w", err)
		}

		authService = service.NewAuth(httpClient, fileStorage)
		secretServise = service.NewSecret(httpClient, fileStorage)

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
