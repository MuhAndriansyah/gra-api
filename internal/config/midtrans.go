package config

import "github.com/spf13/viper"

type MidtransConfig struct {
	ServerKey string
	ClientKey string
	Mode      string
}

func LoadMidtransConfig() MidtransConfig {
	return MidtransConfig{
		ServerKey: viper.GetString("MIDTRANS_SERVER_KEY"),
		ClientKey: viper.GetString("MIDTRANS_CLIENT_KEY"),
		Mode:      viper.GetString("MIDTRANS_MODE"),
	}
}
