package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/julianstephens/license-server/pkg/config"
	"github.com/julianstephens/license-server/pkg/database"
	"github.com/julianstephens/license-server/pkg/logger"
	"github.com/julianstephens/license-server/pkg/redis"
	"github.com/julianstephens/license-server/pkg/router"
	"github.com/spf13/viper"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
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
	if err := config.Setup(); err != nil {
		logger.Fatalf("Could not setup config: %s", err)
	}

	if err := database.Setup(); err != nil {
		logger.Fatalf("could not connect to database: %s", err)
	}

	if err := redis.Init(); err != nil {
		logger.Fatalf("could not connect to redis server: %s", err)
	}

	db := database.GetDB()
	rdb := redis.GetStore()

	r := router.Setup(db, rdb)

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

	logger.Infof("Licensing Server starting at %s:%s", host, port)
	logger.Fatalf("%v", r.Run(host+":"+port))
}
