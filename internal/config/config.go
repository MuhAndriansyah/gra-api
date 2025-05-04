package config

import (
	"fmt"

	"github.com/spf13/viper"
)

type Config struct {
	App      AppConfig
	DB       DBConfig
	Mail     MailConfig
	Redis    RedisConfig
	AWS      AWSConfig
	OAuth    OauthConfig
	Midtrans MidtransConfig
}

func NewConfig(path string) (*Config, error) {
	viper.SetConfigFile(path)
	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
		return nil, fmt.Errorf("error reading config file: %w", err)
	}

	return &Config{
		App:      LoadAppConfig(),
		DB:       LoadDBConfig(),
		Mail:     LoadMailConfig(),
		Redis:    LoadRedisConfig(),
		AWS:      LoadAwsConfig(),
		OAuth:    LoadOauthConfig(),
		Midtrans: LoadMidtransConfig(),
	}, nil
}
