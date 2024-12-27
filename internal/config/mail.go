package config

import "github.com/spf13/viper"

type MailConfig struct {
	SMTPPort     int    `mapstructure:"SMTP_PORT"`
	SMTPHost     string `mapstructure:"SMTP_HOST"`
	MailEmail    string `mapstructure:"MAIL_EMAIL"`
	MailUsername string `mapstructure:"MAIL_USERNAME"`
	MailPassword string `mapstructure:"MAIL_PASSWORD"`
}

func LoadMailConfig() MailConfig {
	return MailConfig{
		SMTPPort:     viper.GetInt("SMTP_PORT"),
		SMTPHost:     viper.GetString("SMTP_HOST"),
		MailEmail:    viper.GetString("MAIL_EMAIL"),
		MailUsername: viper.GetString("MAIL_USERNAME"),
		MailPassword: viper.GetString("MAIL_PASSWORD"),
	}
}
