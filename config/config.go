package config

import (
	"log"
	"os"
	"strconv"
)

type Config struct {
	// Environment Configuration
	ENV        string
	DEBUG_MODE bool

	// Server Configuration
	SERVER_PORT             string
	SERVER_WRITE_TIMEOUT    int // In seconds
	SERVER_READ_TIMEOUT     int // In seconds
	SERVER_IDLE_TIMEOUT     int // In seconds
	SERVER_SHUTDOWN_TIMEOUT int // In seconds

	// Websocket Configuration
	SERVER_WEBSOCKET_PORT             string
	WEBSOCKET_SERVER_SHUTDOWN_TIMEOUT int // In seconds
	WEBSOCKET_URL                     string
	WEBSOCKET_API_KEY                 string

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

	// Auth Configuration
	SEND_VERIFICATION_DELAY_TTL int // In seconds
	MFA_FLAG_TTL                int // In seconds
	ID_TOKEN_HMAC_KEY           string
	ACCESS_TOKEN_TTL            int // In seconds
	REFRESH_TOKEN_TTL           int // In seconds
	MFA_ACCESS_TOKEN_TTL        int // In seconds
	ID_TOKEN_TTL                int // In seconds
	TOTP_PERIOD                 int // In seconds
	SEND_OTP_MAX_RATE_LIMIT     int
	SEND_OTP_MAX_RATE_LIMIT_TTL int // In seconds
	SEND_OTP_DELAY_TTL          int // In seconds

	// Signoz Configuration
	SIGNOZ_URL               string
	SIGNOZ_SERVICE_NAME      string
	SIGNOZ_SERVICE_NAMESPACE string
	SIGNOZ_TRACE_SAMPLE_RATE float64
}

func InitConfig() *Config {
	return &Config{
		ENV:                               os.Getenv("ENV"),
		DEBUG_MODE:                        parseBoolConfig("DEBUG_MODE"),
		SERVER_PORT:                       os.Getenv("SERVER_PORT"),
		SERVER_WRITE_TIMEOUT:              parseIntConfig("SERVER_WRITE_TIMEOUT", 30),
		SERVER_READ_TIMEOUT:               parseIntConfig("SERVER_READ_TIMEOUT", 30),
		SERVER_IDLE_TIMEOUT:               parseIntConfig("SERVER_IDLE_TIMEOUT", 30),
		SERVER_SHUTDOWN_TIMEOUT:           parseIntConfig("SERVER_SHUTDOWN_TIMEOUT", 30),
		SERVER_WEBSOCKET_PORT:             os.Getenv("SERVER_WEBSOCKET_PORT"),
		WEBSOCKET_SERVER_SHUTDOWN_TIMEOUT: parseIntConfig("WEBSOCKET_SERVER_SHUTDOWN_TIMEOUT", 30),
		WEBSOCKET_URL:                     os.Getenv("WEBSOCKET_URL"),
		WEBSOCKET_API_KEY:                 os.Getenv("WEBSOCKET_API_KEY"),
		SMTP_HOST:                         os.Getenv("SMTP_HOST"),
		SMTP_PORT:                         parseIntConfig("SMTP_PORT", 0),
		SMTP_PASSWORD:                     os.Getenv("SMTP_PASSWORD"),
		SMTP_USERNAME:                     os.Getenv("SMTP_USERNAME"),
		SMTP_SENDER:                       os.Getenv("SMTP_SENDER"),
		SMTP_SENDER_NAME:                  os.Getenv("SMTP_SENDER_NAME"),
		DB_USER:                           os.Getenv("DB_USER"),
		DB_PASSWORD:                       os.Getenv("DB_PASSWORD"),
		DB_HOST:                           os.Getenv("DB_HOST"),
		DB_PORT:                           os.Getenv("DB_PORT"),
		DB_NAME:                           os.Getenv("DB_NAME"),
		DB_SSLMODE:                        os.Getenv("DB_SSLMODE"),
		DB_TIMEZONE:                       os.Getenv("DB_TIMEZONE"),
		STORAGE_VENDOR:                    os.Getenv("STORAGE_VENDOR"),
		STORAGE_BUCKET_NAME:               os.Getenv("STORAGE_BUCKET_NAME"),
		STORAGE_PUBLIC_BUCKET_NAME:        os.Getenv("STORAGE_PUBLIC_BUCKET_NAME"),
		ALIYUN_ENDPOINT:                   os.Getenv("ALIYUN_ENDPOINT"),
		ALIYUN_ACCESS_KEY_ID:              os.Getenv("ALIYUN_ACCESS_KEY_ID"),
		ALIYUN_ACCESS_KEY_SECRET:          os.Getenv("ALIYUN_ACCESS_KEY_SECRET"),
		ALIYUN_USER:                       os.Getenv("ALIYUN_USER"),
		MINIO_ENDPOINT:                    os.Getenv("MINIO_ENDPOINT"),
		MINIO_USERNAME:                    os.Getenv("MINIO_USERNAME"),
		MINIO_PASSWORD:                    os.Getenv("MINIO_PASSWORD"),
		MINIO_SSL:                         parseBoolConfig("MINIO_SSL"),
		REDIS_HOST:                        os.Getenv("REDIS_HOST"),
		REDIS_PORT:                        os.Getenv("REDIS_PORT"),
		REDIS_PASSWORD:                    os.Getenv("REDIS_PASSWORD"),
		LOG_PATH:                          os.Getenv("LOG_PATH"),
		SEND_VERIFICATION_DELAY_TTL:       parseIntConfig("SEND_VERIFICATION_DELAY_TTL", 60),
		MFA_FLAG_TTL:                      parseIntConfig("MFA_FLAG_TTL", 604800),
		ID_TOKEN_HMAC_KEY:                 os.Getenv("ID_TOKEN_HMAC_KEY"),
		ACCESS_TOKEN_TTL:                  parseIntConfig("ACCESS_TOKEN_TTL", 86400),
		REFRESH_TOKEN_TTL:                 parseIntConfig("REFRESH_TOKEN_TTL", 604800),
		MFA_ACCESS_TOKEN_TTL:              parseIntConfig("MFA_ACCESS_TOKEN_TTL", 3600),
		ID_TOKEN_TTL:                      parseIntConfig("ID_TOKEN_TTL", 86400),
		TOTP_PERIOD:                       parseIntConfig("TOTP_PERIOD", 120),
		SEND_OTP_MAX_RATE_LIMIT:           parseIntConfig("SEND_OTP_MAX_RATE_LIMIT", 3),
		SEND_OTP_MAX_RATE_LIMIT_TTL:       parseIntConfig("SEND_OTP_MAX_RATE_LIMIT_TTL", 3600),
		SEND_OTP_DELAY_TTL:                parseIntConfig("SEND_OTP_DELAY_TTL", 120),
		SIGNOZ_URL:                        os.Getenv("SIGNOZ_URL"),
		SIGNOZ_SERVICE_NAME:               os.Getenv("SIGNOZ_SERVICE_NAME"),
		SIGNOZ_SERVICE_NAMESPACE:          os.Getenv("SIGNOZ_SERVICE_NAMESPACE"),
		SIGNOZ_TRACE_SAMPLE_RATE:          parseFloatConfig("SIGNOZ_TRACE_SAMPLE_RATE", 120),
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

func parseFloatConfig(envName string, defaultValue float64) float64 {
	envValue := os.Getenv(envName)
	if envValue != "" {
		envValueFloat, err := strconv.ParseFloat(envValue, 64)
		if err != nil {
			log.Fatalf("failed parsing config: %s", envName)
		}
		return envValueFloat
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
