// Пакет конфигурации сервера.
package config

import (
	"flag"
	"fmt"
	"os"
	"time"

	"github.com/ilyakaznacheev/cleanenv"
	"github.com/inhies/go-bytesize"
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
	// Время истечения годности токена
	TokenTTL time.Duration `yaml:"token_ttl" env:"TOKEN_TTL" env-default:"24h"`
	// Максимальный размер тела запроса для регистрации и логина в систему в байтах
	AuthBodyMaxSize int64
}

// Промежуточная конфигурация, служит для преобразования пользовательского ввода типа 10MB
// в реальное занчение конфига
type rawConfig struct {
	AuthBodyMaxSize string `yaml:"auth_body_max_size" env:"AUTH_BODY_MAX_SIZE" env-default:"4KB"`
}

// Загружает конфигурацию из файла, переменных среды и флагов (в порядке приоритета).
func Load() (*Config, error) {
	cfg := &Config{}
	rawCfg := &rawConfig{}

	flag.CommandLine = flag.NewFlagSet("", flag.ContinueOnError)

	var (
		flagC = flag.String("c", "", "config file path")
		flagA = flag.String("a", "", "address to serve https")
		flagD = flag.String("d", "", "database dsn")
		flagT = flag.String("t", "24h", "token ttl in 15h04m05s format")
		flagS = flag.String("s", "4KB", "max auth body size")
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
		err = cleanenv.ReadConfig(configPath, rawCfg)
		if err != nil {
			return nil, fmt.Errorf("failed to parse config file %w", err)
		}
	} else {
		// Применяем переменные окружения
		err := cleanenv.ReadEnv(cfg)
		if err != nil {
			return nil, fmt.Errorf("failed to read environment variables %w", err)
		}
		err = cleanenv.ReadEnv(rawCfg)
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
		case "t":
			d, err := time.ParseDuration(*flagT)
			if err == nil {
				cfg.TokenTTL = d
			}
		case "s":
			rawCfg.AuthBodyMaxSize = *flagS
		}
	})

	// Преобразуем AuthBodyMaxSize из пользовательского ввода в количество байт
	b, err := bytesize.Parse(rawCfg.AuthBodyMaxSize)
	if err != nil {
		return nil, fmt.Errorf("failed to parse AuthBodyMaxSize")
	}
	cfg.AuthBodyMaxSize = int64(b)

	return cfg, nil
}
