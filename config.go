package main

import (
	"database/sql"
	"fmt"
	"github.com/mehmetkule/go-restapi/logger"
	"go.uber.org/zap"
)

// DbConfig is configu struct for database
type DbConfig struct {
	User     string `yaml:"user" envconfig:"DB_USER"`
	Password string `yaml:"password" envconfig:"DB_PASSWORD"`
	DbName   string `yaml:"name" envconfig:"DB_NAME"`
}

// GetDatabase creates database connection using postgres driver
func (c *DbConfig) GetDatabase() (*sql.DB, error) {
	dbInfo := fmt.Sprintf("user=%s password=%s dbname=%s sslmode=disable",
		c.User, c.Password, c.DbName)
	logger.Logger().With(zap.String("dbname",c.DbName)).Warn("Connecting to database.")
	return sql.Open("postgres", dbInfo)
}

