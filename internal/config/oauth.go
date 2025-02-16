package config

import "github.com/spf13/viper"

type OauthConfig struct {
	GoogleClientID     string
	GoogleClientSecret string
	GoogleRedirectUrl  string
}

func LoadOauthConfig() OauthConfig {
	return OauthConfig{
		GoogleClientID:     viper.GetString("GOOGLE_CLIENT_ID"),
		GoogleClientSecret: viper.GetString("GOOGLE_CLIENT_SECRET"),
		GoogleRedirectUrl:  viper.GetString("GOOGLE_REDIRECT_URL"),
	}
}
