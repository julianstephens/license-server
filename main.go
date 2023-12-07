package main

import (
	"fmt"
	"net/http"
	"os"

	"log/slog"

	"github.com/gin-gonic/gin"
	"github.com/julianstephens/license-server/pkg/config"
	"github.com/julianstephens/license-server/pkg/database"
	"github.com/julianstephens/license-server/pkg/logger"
	"github.com/julianstephens/license-server/pkg/router"
	"github.com/spf13/viper"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

var (
	appLogger *slog.Logger
)

//	@title			License Server API
//	@version		0.1.0
//	@description	REST API for managing software licenses

// @host		localhost:8080
// @BasePath	/api/v1
// @schemes http
// @securityDefinitions.apikey ApiKey
// @in header
// @name X-API-KEY
// @descripition User-specific API key
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

	r.GET("/api/v1/docs/*any",
		ginSwagger.WrapHandler(swaggerFiles.Handler),
	)

	r.NoRoute(func(c *gin.Context) {
		c.JSON(http.StatusNotFound, gin.H{"code": "PAGE_NOT_FOUND", "message": "page not found"})
	})

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
