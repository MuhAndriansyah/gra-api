package config

import "github.com/spf13/viper"

type AWSConfig struct {
	AccessKeyID     string
	SecretAccessKey string
	Region          string
	Bucket          string
}

func LoadAwsConfig() AWSConfig {
	return AWSConfig{
		AccessKeyID:     viper.GetString("ACCESS_KEY_ID"),
		SecretAccessKey: viper.GetString("SECRET_ACCESS_KEY"),
		Region:          viper.GetString("AWS_REGION"),
		Bucket:          viper.GetString("AWS_BUCKET"),
	}
}
