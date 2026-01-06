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
	AuthUrl        string `mapstructure:"base_url"`
	TimeoutSeconds int    `mapstructure:"timeout_seconds"`
}

type RedisConfig struct {
	Host            string `mapstructure:"host"`
	Port            string `mapstructure:"port"`
	Password        string `mapstructure:"password"`
	Timeout_seconds int    `mapstructure:"timeout_seconds"`
}

type RedisDbNums struct {
	CurrencyDb int `mapstructure:"currency_db"`
}

type MetricsConfig struct {
	Port string `mapstructure:"port"`
}

type AppConfig struct {
	Service          ServiceConfig        `mapstructure:"gateway_service_config"`
	GrpcClientConfig GrpcClientConfig     `mapstructure:"grpc_client_config"`
	AuthApi          AuthApiServiceConfig `mapstructure:"auth_api_service"`
	RedisConfig      RedisConfig          `mapstructure:"redis_config"`
	RedisDbNums      RedisDbNums          `mapstructure:"Redis_db_nums"`
	MetricsConfig    MetricsConfig        `mapstructure:"metrics_service_config"`
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
