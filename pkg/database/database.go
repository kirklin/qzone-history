package database

import (
	"gorm.io/gorm"
)

type Database interface {
	Connect(config *Config) error
	Close() error
	DB() *gorm.DB
}

type Config struct {
	Host     string
	Port     int
	User     string
	Password string
	DBName   string
	SSLMode  string
}
