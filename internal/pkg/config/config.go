package config

import (
	"log"
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

func Load() {
	AppConfig = Config{
		AppPort:            os.Getenv("APP_PORT"),
		MySQLHost:          os.Getenv("MYSQL_HOST"),
		MySQLPort:          os.Getenv("MYSQL_PORT"),
		MySQLDatabase:      os.Getenv("MYSQL_DATABASE"),
		MySQLUser:          os.Getenv("MYSQL_USER"),
		MySQLPassword:      os.Getenv("MYSQL_PASSWORD"),
		BaseURL:            os.Getenv("BASE_URL"),
		TurnstileSecretKey: os.Getenv("TURNSTILE_SECRET_KEY"),
	}

	required := map[string]string{
		"APP_PORT":             AppConfig.AppPort,
		"MYSQL_HOST":           AppConfig.MySQLHost,
		"MYSQL_PORT":           AppConfig.MySQLPort,
		"MYSQL_DATABASE":       AppConfig.MySQLDatabase,
		"MYSQL_USER":           AppConfig.MySQLUser,
		"MYSQL_PASSWORD":       AppConfig.MySQLPassword,
		"BASE_URL":             AppConfig.BaseURL,
		"TURNSTILE_SECRET_KEY": AppConfig.TurnstileSecretKey,
	}

	var missing []string
	for _, key := range []string{
		"APP_PORT", "MYSQL_HOST", "MYSQL_PORT", "MYSQL_DATABASE",
		"MYSQL_USER", "MYSQL_PASSWORD", "BASE_URL", "TURNSTILE_SECRET_KEY",
	} {
		if required[key] == "" {
			missing = append(missing, key)
		}
	}

	if len(missing) > 0 {
		log.Fatalf("missing required environment variables: %v", missing)
	}
}
