package viper

import (
	"fmt"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
	"log"
)

// RunViperExample tiny example how to use viper package to handle flags, envs and config files
func RunViperExample() {
	//flag definitions, --help is automatically added
	pflag.StringP("username", "u", "", "The username of the user")
	pflag.StringP("password", "p", "", "The password of the user")

	//actual parsing from args [1:]
	pflag.Parse()

	//binding parsed flags to viper package
	err := viper.BindPFlags(pflag.CommandLine)

	if err != nil {
		log.Fatalf("error binding flags: %v", err)
	}

	//now can be easily accessed from viper
	fmt.Printf("provided username: %v\n", viper.GetString("username"))
	fmt.Printf("provided password: %v\n", viper.GetString("password"))

	//viper can also bind ENVs
	viper.MustBindEnv("PATH")
	fmt.Printf("PATH: %v\n", viper.Get("PATH"))

	//viper can also read config files
	viper.SetConfigType("yaml")
	viper.SetConfigFile("cmd/viper/config.yaml")
	err = viper.ReadInConfig()
	if err != nil {
		log.Fatalf("error reading config: %v", err)
	}
	var config PostgresConfig
	err = viper.UnmarshalKey("db.postgres", &config)
	if err != nil {
		log.Fatalf("error unmarshalling config: %v", err)
	}
	fmt.Printf("config: %v\n", config)
}

type PostgresConfig struct {
	Host     string `mapstructure:"host"`
	Port     int    `mapstructure:"port"`
	User     string `mapstructure:"user"`
	Password string `mapstructure:"password"`
}
