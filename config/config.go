package config

import "os"

type Config struct {
	// Server Configuration
	SERVER_PORT string

	// Database Configuration
	DB_USER     string
	DB_PASSWORD string
	DB_HOST     string
	DB_PORT     string
	DB_NAME     string
	DB_SSLMODE  string
	DB_TIMEZONE string
}

func InitConfig() *Config {
	return &Config{
		SERVER_PORT: os.Getenv("SERVER_PORT"),
		DB_USER:     os.Getenv("DB_USER"),
		DB_PASSWORD: os.Getenv("DB_PASSWORD"),
		DB_HOST:     os.Getenv("DB_HOST"),
		DB_PORT:     os.Getenv("DB_PORT"),
		DB_NAME:     os.Getenv("DB_NAME"),
		DB_SSLMODE:  os.Getenv("DB_SSLMODE"),
		DB_TIMEZONE: os.Getenv("DB_TIMEZONE"),
	}
}
