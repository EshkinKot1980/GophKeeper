// Пакет cli содержит интерфейс омандной строки.
package cli

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/EshkinKot1980/GophKeeper/internal/client/cli/utils"
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
	Upload(secret dto.SecretRequest, data []byte) error
	// GetSecretAndInfo получает секрет пользователя с сервера по id,
	// возвращает расшиврованные данные в виде []byte и информацию о секрете
	GetSecretAndInfo(id uint64) ([]byte, dto.SecretInfo, error)
	// InfoList получает информацию о всех секретах пользователя с сервера.
	InfoList() ([]dto.SecretInfo, error)
}

// Prompt обслуживает пользовательский ввод
type Prompt interface {
	// SecretName ввод названия секрета
	SecretName() (string, error)
	// RegisterCredentials ввод учетных данных для регистрации
	RegisterCredentials() (dto.Credentials, error)
	// Credentials ввод учетных данных для входа или сохранения в системе
	Credentials() (dto.Credentials, error)
	// Overwrite() запрашивает у пользователя нужно ли файл переписать
	Overwrite(fileName string) bool
}

var (
	cfg           *config.Config
	authService   AuthService
	secretService SecretService
	prompt        Prompt
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
		secretService = service.NewSecret(httpClient, fileStorage)
		prompt = utils.NewPrompt()

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
