package config

import (
	"fmt"
	"log"

	"github.com/spf13/viper"
)

type ServiceConfig struct {
	Host           string `mapstructure:"host"`
	Port           string `mapstructure:"port"`
	TimeoutSeconds int    `mapstructure:"timeout_seconds"`
}

type GrpcClientConfig struct {
	CurrencyServiceUrl string `mapstructure:"currency_service_url"`
	TimeoutSeconds     int    `mapstructure:"timeout_seconds"`
}

type AuthApiServiceConfig struct {
	AuthUrl string `mapstructure:"base_url"`
}

type AppConfig struct {
	Service          ServiceConfig        `mapstructure:"gateway_service_config"`
	GrpcClientConfig GrpcClientConfig     `mapstructure:"grpc_client_config"`
	AuthApi          AuthApiServiceConfig `mapstructure:"auth_api_service"`
}

func LoadConfig(path string) (*AppConfig, error) {
	var config AppConfig

	viper.SetConfigFile(path)

	// Reading the configuration file
	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			log.Println("Config file not found")
		} else {
			return &config, fmt.Errorf("error reading config file: %w", err)
		}
	}

	// Unmarshal into structure
	if err := viper.Unmarshal(&config); err != nil {
		return &config, fmt.Errorf("unable to unmarshal config: %w", err)
	}

	return &config, nil
}
