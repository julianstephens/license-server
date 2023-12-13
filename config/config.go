package config

import (
	"strings"

	"github.com/julianstephens/license-server/internal/model"
	"github.com/spf13/viper"
)

var Config *model.Config

func Setup() error {
	var conf model.Config

	v := viper.New()
	v.AutomaticEnv()
	v.SetEnvPrefix("ls")
	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	v.AddConfigPath(".")
	v.SetConfigName("config")
	v.SetConfigType("yaml")

	if err := v.ReadInConfig(); err != nil {
		return err
	}

	if err := v.Unmarshal(&conf); err != nil {
		return err
	}

	Config = &conf

	return nil
}

func GetConfig() *model.Config {
	return Config
}
