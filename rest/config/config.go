package config

import (
	"fmt"
	"github.com/spf13/viper"
	"log"
	"time"
)

type AppConfig struct {
	Server ServerConfig `mapstructure:"server"`
	DB     DBConfig     `mapstructure:"db"`
}

type ServerConfig struct {
	Host string `mapstructure:"host"`
	Port int    `mapstructure:"port"`
}

type DBConfig struct {
	User     string        `mapstructure:"user"`
	Password string        `mapstructure:"password"`
	Host     string        `mapstructure:"host"`
	Port     int           `mapstructure:"port"`
	Database string        `mapstructure:"database"`
	Timeout  time.Duration `mapstructure:"timeout"`
	PoolMax  int           `mapstructure:"pool_max_conns"`
	PoolMin  int           `mapstructure:"pool_min_conns"`
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
