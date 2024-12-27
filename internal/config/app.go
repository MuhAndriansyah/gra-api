package config

import (
	"github.com/spf13/viper"
)

type AppConfig struct {
	Name          string
	Env           string
	LogLevel      string
	JWTPrivateKey string
}

func LoadAppConfig() AppConfig {
	return AppConfig{
		Name:          viper.GetString("APP_NAME"),
		Env:           viper.GetString("APP_ENV"),
		LogLevel:      viper.GetString("LOG_LEVEL"),
		JWTPrivateKey: viper.GetString("JWT_PRIVATE_KEY"),
	}
}
