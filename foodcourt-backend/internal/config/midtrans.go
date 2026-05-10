package config

import (
	"os"
)

type MidtransConfig struct {
	ServerKey    string
	ClientKey    string
	IsProduction bool
}

func LoadMidtransConfig() *MidtransConfig {
	return &MidtransConfig{
		ServerKey:    os.Getenv("MIDTRANS_SERVER_KEY"),
		ClientKey:    os.Getenv("MIDTRANS_CLIENT_KEY"),
		IsProduction: os.Getenv("MIDTRANS_IS_PRODUCTION") == "true",
	}
}
