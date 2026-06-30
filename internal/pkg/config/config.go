package config

import (
	"fmt"
	"log/slog"
	"os"
)

type Config struct {
	AppPort            string
	MySQLHost          string
	MySQLPort          string
	MySQLDatabase      string
	MySQLUser          string
	MySQLPassword      string
	BaseURL            string
	TurnstileSecretKey string
}

var AppConfig Config

func Parse() (Config, error) {
	cfg := Config{
		AppPort:            os.Getenv("APP_PORT"),
		MySQLHost:          os.Getenv("MYSQL_HOST"),
		MySQLPort:          os.Getenv("MYSQL_PORT"),
		MySQLDatabase:      os.Getenv("MYSQL_DATABASE"),
		MySQLUser:          os.Getenv("MYSQL_USER"),
		MySQLPassword:      os.Getenv("MYSQL_PASSWORD"),
		BaseURL:            os.Getenv("BASE_URL"),
		TurnstileSecretKey: os.Getenv("TURNSTILE_SECRET_KEY"),
	}

	values := map[string]string{
		"APP_PORT":             cfg.AppPort,
		"MYSQL_HOST":           cfg.MySQLHost,
		"MYSQL_PORT":           cfg.MySQLPort,
		"MYSQL_DATABASE":       cfg.MySQLDatabase,
		"MYSQL_USER":           cfg.MySQLUser,
		"MYSQL_PASSWORD":       cfg.MySQLPassword,
		"BASE_URL":             cfg.BaseURL,
		"TURNSTILE_SECRET_KEY": cfg.TurnstileSecretKey,
	}

	var missing []string
	for _, key := range []string{
		"APP_PORT", "MYSQL_HOST", "MYSQL_PORT", "MYSQL_DATABASE",
		"MYSQL_USER", "MYSQL_PASSWORD", "BASE_URL", "TURNSTILE_SECRET_KEY",
	} {
		if values[key] == "" {
			missing = append(missing, key)
		}
	}

	if len(missing) > 0 {
		return Config{}, fmt.Errorf("missing required environment variables: %v", missing)
	}

	return cfg, nil
}

func Load() {
	cfg, err := Parse()
	if err != nil {
		slog.Error("failed to load config", "error", err)
		os.Exit(1)
	}
	AppConfig = cfg
}
