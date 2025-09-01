package config

import (
	"log"
	"os"
	"strconv"
)

type Config struct {
	// Environment Configuration
	ENV string

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

	// Storage Configuration
	STORAGE_VENDOR             string
	STORAGE_BUCKET_NAME        string
	STORAGE_PUBLIC_BUCKET_NAME string
	// Aliyun
	ALIYUN_ENDPOINT          string
	ALIYUN_ACCESS_KEY_ID     string
	ALIYUN_ACCESS_KEY_SECRET string
	ALIYUN_USER              string
	// Minio
	MINIO_ENDPOINT string
	MINIO_USERNAME string
	MINIO_PASSWORD string
	MINIO_SSL      bool

	// Redis Configuration
	REDIS_HOST     string
	REDIS_PORT     string
	REDIS_PASSWORD string

	// Logging Configuration
	LOG_PATH string
}

func InitConfig() *Config {
	serverWriteTimeout := parseIntConfig("SERVER_WRITE_TIMEOUT", 30)
	serverReadTimeout := parseIntConfig("SERVER_READ_TIMEOUT", 30)
	serverIdleTimeout := parseIntConfig("SERVER_IDLE_TIMEOUT", 30)
	serverShutdownTimeout := parseIntConfig("SERVER_SHUTDOWN_TIMEOUT", 30)
	smtpPort := parseIntConfig("SMTP_PORT", 0)
	minioSsl := parseBoolConfig("MINIO_SSL")

	return &Config{
		ENV:                        os.Getenv("ENV"),
		SERVER_PORT:                os.Getenv("SERVER_PORT"),
		SERVER_WRITE_TIMEOUT:       serverWriteTimeout,
		SERVER_READ_TIMEOUT:        serverReadTimeout,
		SERVER_IDLE_TIMEOUT:        serverIdleTimeout,
		SERVER_SHUTDOWN_TIMEOUT:    serverShutdownTimeout,
		SMTP_HOST:                  os.Getenv("SMTP_HOST"),
		SMTP_PORT:                  smtpPort,
		SMTP_PASSWORD:              os.Getenv("SMTP_PASSWORD"),
		SMTP_USERNAME:              os.Getenv("SMTP_USERNAME"),
		SMTP_SENDER:                os.Getenv("SMTP_SENDER"),
		SMTP_SENDER_NAME:           os.Getenv("SMTP_SENDER_NAME"),
		DB_USER:                    os.Getenv("DB_USER"),
		DB_PASSWORD:                os.Getenv("DB_PASSWORD"),
		DB_HOST:                    os.Getenv("DB_HOST"),
		DB_PORT:                    os.Getenv("DB_PORT"),
		DB_NAME:                    os.Getenv("DB_NAME"),
		DB_SSLMODE:                 os.Getenv("DB_SSLMODE"),
		DB_TIMEZONE:                os.Getenv("DB_TIMEZONE"),
		STORAGE_VENDOR:             os.Getenv("STORAGE_VENDOR"),
		STORAGE_BUCKET_NAME:        os.Getenv("STORAGE_BUCKET_NAME"),
		STORAGE_PUBLIC_BUCKET_NAME: os.Getenv("STORAGE_PUBLIC_BUCKET_NAME"),
		ALIYUN_ENDPOINT:            os.Getenv("ALIYUN_ENDPOINT"),
		ALIYUN_ACCESS_KEY_ID:       os.Getenv("ALIYUN_ACCESS_KEY_ID"),
		ALIYUN_ACCESS_KEY_SECRET:   os.Getenv("ALIYUN_ACCESS_KEY_SECRET"),
		ALIYUN_USER:                os.Getenv("ALIYUN_USER"),
		MINIO_ENDPOINT:             os.Getenv("MINIO_ENDPOINT"),
		MINIO_USERNAME:             os.Getenv("MINIO_USERNAME"),
		MINIO_PASSWORD:             os.Getenv("MINIO_PASSWORD"),
		MINIO_SSL:                  minioSsl,
		REDIS_HOST:                 os.Getenv("REDIS_HOST"),
		REDIS_PORT:                 os.Getenv("REDIS_PORT"),
		REDIS_PASSWORD:             os.Getenv("REDIS_PASSWORD"),
		LOG_PATH:                   os.Getenv("LOG_PATH"),
	}
}

func parseIntConfig(envName string, defaultValue int) int {
	envValue := os.Getenv(envName)
	if envValue != "" {
		envValueInt, err := strconv.Atoi(envValue)
		if err != nil {
			log.Fatalf("failed parsing config: %s", envName)
		}
		return envValueInt
	}
	return defaultValue
}

func parseBoolConfig(envName string) bool {
	envValue := os.Getenv(envName)
	if envValue != "" {
		envValueBool, err := strconv.ParseBool(envValue)
		if err != nil {
			log.Fatalf("failed parsing config: %s", envName)
		}
		return envValueBool
	}
	return false
}
