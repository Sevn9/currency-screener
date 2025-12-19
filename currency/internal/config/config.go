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

type ManagementConfig struct {
	Host           string `mapstructure:"host"`
	Port           string `mapstructure:"port"`
	TimeoutSeconds int    `mapstructure:"timeout_seconds"`
}

type MetricsConfig struct {
	Port string `mapstructure:"port"`
}

type PublicCurrencyAPIConfig struct {
	ApiUrl         string `mapstructure:"api_url"`
	TimeoutSeconds int    `mapstructure:"timeout_seconds"`
}

type DatabasePostgresConfig struct {
	Host     string `mapstructure:"host"`
	Port     int    `mapstructure:"port"`
	Username string `mapstructure:"username"`
	Password string `mapstructure:"password"`
	DBName   string `mapstructure:"dbname"`
}

type WorkerConfig struct {
	Schedule     string `mapstructure:"schedule"`
	Timeout      int    `mapstructure:"timeout"`
	CurrencyPair struct {
		BaseCurrency   string `mapstructure:"base_currency"`
		TargetCurrency string `mapstructure:"target_currency"`
	} `mapstructure:"currency_pair"`
}

// type LoggingConfig struct {
// 	Level  string `mapstructure:"level"`
// 	Format string `mapstructure:"format"`
// }

type AppConfig struct {
	Service           ServiceConfig           `mapstructure:"grpc_service_config"`
	PublicCurrencyApi PublicCurrencyAPIConfig `mapstructure:"currency_api_config"`
	ManagementService ManagementConfig        `mapstructure:"management_service_config"`
	PostgresDb        DatabasePostgresConfig  `mapstructure:"database_postgres"`
	Worker            WorkerConfig            `mapstructure:"worker"`
	MetricsConfig     MetricsConfig           `mapstructure:"metrics_service_config"`
}

func LoadConfig(path string) (*AppConfig, error) {
	var config AppConfig

	viper.SetConfigFile(path)

	// Читаем конфигурационный файл
	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			log.Println("Config file not found")
		} else {
			return &config, fmt.Errorf("error reading config file: %w", err)
		}
	}

	// Unmarshal в структуру
	if err := viper.Unmarshal(&config); err != nil {
		return &config, fmt.Errorf("unable to unmarshal config: %w", err)
	}

	return &config, nil
}
