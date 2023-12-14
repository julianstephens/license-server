package config

import (
	"strings"
	"sync"

	"github.com/julianstephens/license-server/internal/model"
	"github.com/julianstephens/license-server/pkg/logger"
	"github.com/spf13/viper"
)

var (
	Config *model.Config
	once   sync.Once
	err    error
)

func GetConfig() *model.Config {
	once.Do(func() {
		Config, err = setup()
		if err != nil {
			logger.Fatalf("unable to load application environment: %+v", err)
		}
	})

	return Config
}

func setup() (*model.Config, error) {
	var conf model.Config

	v := viper.New()
	v.AutomaticEnv()
	v.SetEnvPrefix("ls")
	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	v.AddConfigPath(".")
	v.SetConfigName("config")
	v.SetConfigType("yaml")

	if err = v.ReadInConfig(); err != nil {
		return nil, err
	}

	if err := v.Unmarshal(&conf); err != nil {
		return nil, err
	}

	return &conf, nil
}
