package config

import (
	"fmt"

	"github.com/spf13/viper"
)

type DBConfig struct {
	DBUser       string
	DBPassword   string
	DBHost       string
	DBPort       string
	DBName       string
	MaxIdleConns int
	MaxOpenConns int
	DSN          string
}

func LoadDBConfig() DBConfig {
	dsn := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable",
		viper.GetString("DB_PGX_USER"),
		viper.GetString("DB_PGX_PASSWORD"),
		viper.GetString("DB_PGX_HOST"),
		viper.GetString("DB_PGX_PORT"),
		viper.GetString("DB_PGX_NAME"),
	)

	return DBConfig{
		DBUser:       viper.GetString("DB_PGX_USER"),
		DBPassword:   viper.GetString("DB_PGX_PASSWORD"),
		DBHost:       viper.GetString("DB_PGX_HOST"),
		DBPort:       viper.GetString("DB_PGX_PORT"),
		DBName:       viper.GetString("DB_PGX_NAME"),
		MaxIdleConns: viper.GetInt("DB_MAX_IDLE_CONNS"),
		MaxOpenConns: viper.GetInt("DB_MAX_OPEN_CONNS"),
		DSN:          dsn,
	}
}
