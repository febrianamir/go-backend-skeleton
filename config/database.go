package config

import (
	"app/lib"
	"fmt"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func (c *Config) NewDB() (*lib.Database, error) {
	dsn := c.getPostgresDSN()
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	return &lib.Database{DB: db}, nil
}

func (c *Config) getPostgresDSN() string {
	return fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=%s TimeZone=%s",
		c.DB_HOST, c.DB_USER, c.DB_PASSWORD, c.DB_NAME, c.DB_PORT, c.DB_SSLMODE, c.DB_TIMEZONE)
}
