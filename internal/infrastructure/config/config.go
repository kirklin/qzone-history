package config

import (
	"github.com/spf13/viper"
)

type Config struct {
	Database struct {
		Type     string // 数据库类型，例如 "sqlite", "postgres", "mysql"
		Host     string
		Port     int
		User     string
		Password string
		DBName   string
		SSLMode  string
	}
	QzoneAPI struct {
		QRCodeURL  string
		LoginURL   string
		MessageURL string
	}
}

func LoadConfig() (*Config, error) {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath("./config")

	err := viper.ReadInConfig()
	if err != nil {
		return nil, err
	}

	var config Config
	err = viper.Unmarshal(&config)
	if err != nil {
		return nil, err
	}

	return &config, nil
}
