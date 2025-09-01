package logger

import (
	"time"

	"go.uber.org/zap"
)

const (
	LOGGER_ENV_SETUP_DEVELOPMENT_VALUE     = "development"
	LOGGER_ENV_SETUP_NON_DEVELOPMENT_VALUE = "non_development"
)

type messageLog struct {
	Timestamps time.Time   `json:"timestamps"`
	Level      string      `json:"level"`
	Message    string      `json:"message"`
	Fields     []zap.Field `json:"fields"`
}

type LoggerSetup struct {
	Env  string `json:"env"`
	Path string `json:"path"`
}

func (setup *LoggerSetup) valueDefault() {
	if setup.Path == "" {
		setup.Path = "logs/"
	}

	if setup.Env == "" {
		setup.Env = LOGGER_ENV_SETUP_DEVELOPMENT_VALUE
	}
}
