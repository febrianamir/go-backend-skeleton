package config

import (
	"log"
	"os"
	"strconv"
)

type Config struct {
	// Server Configuration
	SERVER_PORT             string
	SERVER_WRITE_TIMEOUT    int // In seconds
	SERVER_READ_TIMEOUT     int // In seconds
	SERVER_IDLE_TIMEOUT     int // In seconds
	SERVER_SHUTDOWN_TIMEOUT int // In seconds

	// Database Configuration
	DB_USER     string
	DB_PASSWORD string
	DB_HOST     string
	DB_PORT     string
	DB_NAME     string
	DB_SSLMODE  string
	DB_TIMEZONE string

	// SMTP Configuration
	SMTP_HOST        string
	SMTP_PORT        int
	SMTP_PASSWORD    string
	SMTP_USERNAME    string
	SMTP_SENDER      string
	SMTP_SENDER_NAME string
}

func InitConfig() *Config {
	serverWriteTimeout := parseIntConfig("SERVER_WRITE_TIMEOUT", 30)
	serverReadTimeout := parseIntConfig("SERVER_READ_TIMEOUT", 30)
	serverIdleTimeout := parseIntConfig("SERVER_IDLE_TIMEOUT", 30)
	serverShutdownTimeout := parseIntConfig("SERVER_SHUTDOWN_TIMEOUT", 30)
	smtpPort := parseIntConfig("SMTP_PORT", 0)

	return &Config{
		SERVER_PORT:             os.Getenv("SERVER_PORT"),
		SERVER_WRITE_TIMEOUT:    serverWriteTimeout,
		SERVER_READ_TIMEOUT:     serverReadTimeout,
		SERVER_IDLE_TIMEOUT:     serverIdleTimeout,
		SERVER_SHUTDOWN_TIMEOUT: serverShutdownTimeout,
		SMTP_HOST:               os.Getenv("SMTP_HOST"),
		SMTP_PORT:               smtpPort,
		SMTP_PASSWORD:           os.Getenv("SMTP_PASSWORD"),
		SMTP_USERNAME:           os.Getenv("SMTP_USERNAME"),
		SMTP_SENDER:             os.Getenv("SMTP_SENDER"),
		SMTP_SENDER_NAME:        os.Getenv("SMTP_SENDER_NAME"),
		DB_USER:                 os.Getenv("DB_USER"),
		DB_PASSWORD:             os.Getenv("DB_PASSWORD"),
		DB_HOST:                 os.Getenv("DB_HOST"),
		DB_PORT:                 os.Getenv("DB_PORT"),
		DB_NAME:                 os.Getenv("DB_NAME"),
		DB_SSLMODE:              os.Getenv("DB_SSLMODE"),
		DB_TIMEZONE:             os.Getenv("DB_TIMEZONE"),
	}
}

func parseIntConfig(envName string, defaultValue int) int {
	envValue := os.Getenv(envName)
	if envValue != "" {
		envValueInt, err := strconv.Atoi(envValue)
		if err != nil {
			log.Fatal("failed parsing config: SERVER_SHUTDOWN_TIMEOUT")
		}
		return envValueInt
	}
	return defaultValue
}
