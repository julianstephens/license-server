package config

import (
	"log"
	"strings"
	"sync"

	"github.com/spf13/viper"

	"github.com/julianstephens/license-server/backend/pkg/logger"
	"github.com/julianstephens/license-server/backend/pkg/model"
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

	v := findConfigFile()

	if err := v.Unmarshal(&conf); err != nil {
		return nil, err
	}

	return &conf, nil
}

func findConfigFile() *viper.Viper {
	var err error
	path := "./"
	v := viper.New()
	for i := 0; i < 10; i++ {

		v.AddConfigPath(path)
		v.SetConfigName("config")
		v.SetConfigType("yaml")
		err = v.ReadInConfig()
		if err != nil {
			if strings.Contains(err.Error(), "Not Found") {
				path = path + "../"
				continue
			}
			log.Fatal("panic in config parser : " + err.Error())
		} else {
			break
		}
	}

	return v
}
