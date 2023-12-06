package config

import (
	"log"

	"github.com/spf13/viper"
)

type Config struct {
	AppEnv string `mapstructure:"APP_ENV"`
	DBHost string `mapstructure:"POSTGRES_HOST"`
	DBUser string `mapstructure:"POSTGRES_USER"`
	DBPass string `mapstructure:"POSTGRES_PASSWORD"`
	DBName string `mapstructure:"POSTGRES_DB"`
}

func LoadConfig() *Config {
	conf := Config{}
	viper.SetConfigFile(".env")

	err := viper.ReadInConfig()
	if err != nil {
		log.Fatal("Can't find the file .env: ", err)
	}

	err = viper.Unmarshal(&conf)
	if err != nil {
		log.Fatal("Environment can't be loaded: ", err)
	}

	if conf.AppEnv == "development" {
		log.Println("Running in development mode")
	}

	return &conf
}
