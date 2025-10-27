package config

import (
	"fmt"
	"log"

	"github.com/spf13/viper"
)

type ServiceConfig struct {
	Host    string `mapstructure:"host"`
	Port    string `mapstructure:"port"`
	Timeout int    `mapstructure:"timeout"`
}

type PublicCurrencyAPI struct {
	ApiUrl string `mapstructure:"api_url"`
}

// type DatabaseConfig struct {
// 	Host     string `mapstructure:"host"`
// 	Port     int    `mapstructure:"port"`
// 	Username string `mapstructure:"username"`
// 	Password string `mapstructure:"password"`
// 	DBName   string `mapstructure:"dbname"`
// }

// type LoggingConfig struct {
// 	Level  string `mapstructure:"level"`
// 	Format string `mapstructure:"format"`
// }

type AppConfig struct {
	Service           ServiceConfig     `mapstructure:"service_config"`
	PublicCurrencyApi PublicCurrencyAPI `mapstructure:"currency_api_config"`
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
