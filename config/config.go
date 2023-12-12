package config

import "github.com/spf13/viper"

func Setup() error {
	viper.SetEnvPrefix("ls")
	viper.AutomaticEnv()
	return nil
}
