// Пакет config конфигурации сервера.
package config

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/ilyakaznacheev/cleanenv"
)

const (
	configDirName  = ".gophkeeper"
	configFileName = "config.yml"
)

// Конфигурация сервера.
type Config struct {
	// Адрес https cервера в формате "host:port".
	ServerAddr string `yaml:"https_addr" env:"SERVER_ADDR" env-default:"localhost:8443"`
	// Разрешить ли самодписанные сертификаты для https соединения с сервером
	AllowSelfSignedCert bool `yaml:"allow_self_signed_cert" env:"ALLOW_SELF_SIGNED_CERT" env-default:"true"`
}

// Load pагружает конфигурацию из файла и переменных.
func Load() (*Config, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return nil, fmt.Errorf("failed to get user home dir: %w", err)
	}

	configDir := filepath.Join(homeDir, configDirName)
	// Создаем директорию конфига, если её нет
	if err := os.MkdirAll(configDir, 0700); err != nil {
		return nil, fmt.Errorf("failed to create config dir: %w", err)
	}

	cfg := &Config{}
	configPath := filepath.Join(configDir, configFileName)

	_, err = os.Stat(configPath)
	if err != nil {
		if !os.IsNotExist(err) {
			return nil, fmt.Errorf("invalid config file path, %w", err)
		}

		// Если файла нет, применяем переменные окружения
		err := cleanenv.ReadEnv(cfg)
		if err != nil {
			return nil, fmt.Errorf("failed to read environment variables %w", err)
		}
	} else {
		// Разбираем файл конфигурации и применяем переменные окружения
		err := cleanenv.ReadConfig(configPath, cfg)
		if err != nil {
			return nil, fmt.Errorf("failed to parse config file %w", err)
		}
	}

	return cfg, nil
}
