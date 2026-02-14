// Пакет конфигурации сервера.
package config

import (
	"flag"
	"fmt"
	"os"

	"github.com/ilyakaznacheev/cleanenv"
)

// Конфигурация сервера.
type Config struct {
	// DSN для подключения к СУБД Postgres.
	DatabaseDSN string `yaml:"database_dsn" env:"DATABASE_DSN"`
	// Адрес для работы https cервера в формате "host:port".
	HTTPSaddr string `yaml:"https_addr" env:"HTTPS_ADDR" env-default:"localhost:8443"`
	// Путь к сертификату для обеспечения `HTTPS/TLS` соединения.
	// Если сертификат подписан центром сертификации,
	// файл сертификата должен представлять собой конкатенацию сертификата сервера,
	// любых промежуточных сертификатов и сертификата центра сертификации.
	TLScert string `yaml:"tls_cert" env:"TLS_CERT" env-default:"rsa/tls.crt"`
	// Путь к приватному ключу для обеспечения `HTTPS/TLS` соединения.
	TLSkey string `yaml:"tls_key" env:"TLS_KEY" env-default:"rsa/tls.key"`
	// Путь к приватному ключу для JWT
	JWTpriv string `yaml:"jwt_priv" env:"JWT_PRIV" env-default:"rsa/jwt-priv.pem"`
	// Путь к публичному ключу для JWT
	JWTpub string `yaml:"jwt_pub" env:"JWT_PUB" env-default:"rsa/jwt-pub.pem"`
}

// Загружает конфигурацию из файла, переменных среды и флагов (в порядке приоритета).
func Load() (*Config, error) {
	cfg := &Config{}

	flag.CommandLine = flag.NewFlagSet("", flag.ContinueOnError)

	var (
		flagC = flag.String("c", "", "config file path")
		flagA = flag.String("a", "", "address to serve https")
		flagD = flag.String("d", "", "database dsn")
	)

	err := flag.CommandLine.Parse(os.Args[1:])
	if err != nil {
		return nil, fmt.Errorf("failed to parse flags %w", err)
	}

	configPath := *flagC
	if configPath != "" {
		// Разбираем файл конфигурации и применяем переменные окружения
		err := cleanenv.ReadConfig(configPath, cfg)
		if err != nil {
			return nil, fmt.Errorf("failed to parse config file %w", err)
		}
	} else {
		// Применяем переменные окружения
		err := cleanenv.ReadEnv(cfg)
		if err != nil {
			return nil, fmt.Errorf("failed to read environment variables %w", err)
		}
	}

	// Применяем заданные флаги
	flag.Visit(func(fl *flag.Flag) {
		switch fl.Name {
		case "a":
			cfg.HTTPSaddr = *flagA
		case "d":
			cfg.DatabaseDSN = *flagD
		}
	})

	return cfg, nil
}
