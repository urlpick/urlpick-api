package config

import "os"

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
}
