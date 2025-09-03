package config

import (
	"app/lib/logger"

	gormLogger "gorm.io/gorm/logger"
)

func (c *Config) NewSQLLogger() *logger.SQLLogger {
	return &logger.SQLLogger{
		Interface: gormLogger.Default.LogMode(gormLogger.Silent),
		Env:       c.ENV,
		DebugMode: c.DEBUG_MODE,
	}
}
