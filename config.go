package main

import (
	"database/sql"
	"fmt"
	"github.com/mehmetkule/go-restapi/logger"
	"go.uber.org/zap"
)

// DbConfig is configu struct for database
type DbConfig struct {
	Host     string `yaml:"host" envconfig:"DB_HOST"`
	Port     string `yaml:"port" envconfig:"DB_PORT"`
	User     string `yaml:"user" envconfig:"DB_USER"`
	Password string `yaml:"password" envconfig:"DB_PASSWORD"`
	DbName   string `yaml:"name" envconfig:"DB_NAME"`
}

// GetDatabase creates database connection using postgres driver
func (c *DbConfig) GetDatabase() (*sql.DB, error) {
	dbInfo := fmt.Sprintf("host=%v port=%v user=%s password=%s dbname=%s sslmode=disable",c.Host,
		c.Port,c.User, c.Password, c.DbName)
	logger.Logger().With(zap.String("dbname",c.DbName)).Warn("Connecting to database.")
	return sql.Open("postgres", dbInfo)
}

