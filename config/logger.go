package config

import (
	"app/lib/logger"
)

func (c *Config) NewLogger() {
	var setup = logger.LoggerSetup{
		Path: c.LOG_PATH,
		Env:  c.ENV,
	}
	logger.Init(setup)
}
