package main

import (
	"fmt"
	"os"

	"log/slog"

	"github.com/julianstephens/license-server/pkg/config"
	"github.com/julianstephens/license-server/pkg/database"
	"github.com/julianstephens/license-server/pkg/logger"
	"github.com/julianstephens/license-server/pkg/router"
	"github.com/spf13/viper"
)

var (
	appLogger *slog.Logger
)

func main() {
	logger.Setup()
	appLogger = logger.GetLogger()

	err := config.Setup()
	if err != nil {
		appLogger.Error("Could not setup config", err)
		os.Exit(1)
	}

	err = database.Setup()
	if err != nil {
		appLogger.Error("could not connect to database", err)
		os.Exit(1)
	}
	db := database.GetDB()

	r := router.Setup(db)

	host := "0.0.0.0"
	port := "8080"

	if h := viper.GetString("server_host"); h != "" {
		host = h
	}
	if p := viper.GetString("server_port"); p != "" {
		port = p
	}

	appLogger.Info(fmt.Sprintf("Licensing Server starting at %s:%s", host, port))
	err = r.Run(host + ":" + port)
	if err != nil {
		appLogger.Error("Could not start server", err)
		os.Exit(1)
	}
}
