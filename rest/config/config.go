package config

import (
	"fmt"
	"github.com/spf13/viper"
	"log"
)

type AppConfig struct {
	Server struct {
		Host string `mapstructure:"host"`
		Port int    `mapstructure:"port"`
	} `mapstructure:"server"`
	DB struct {
		User     string `mapstructure:"user"`
		Password string `mapstructure:"password"`
		Host     string `mapstructure:"host"`
		Port     int    `mapstructure:"port"`
		Database string `mapstructure:"database"`
	} `mapstructure:"db"`
}

func Read(env string) *AppConfig {
	viper.SetConfigType("yaml")
	viper.SetConfigFile(fmt.Sprintf("rest/config-%s.yaml", env))
	err := viper.ReadInConfig()
	if err != nil {
		log.Fatalf("error reading config: %v", err)
	}
	config := new(AppConfig)
	err = viper.Unmarshal(config)
	if err != nil {
		log.Fatalf("error unmarshalling config: %v", err)
	}
	return config
}
