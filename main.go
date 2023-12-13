package main

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/julianstephens/license-server/config"
	"github.com/julianstephens/license-server/internal/router"
	"github.com/julianstephens/license-server/pkg/database"
	"github.com/julianstephens/license-server/pkg/logger"
	"github.com/julianstephens/license-server/pkg/redis"
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
	conf := config.GetConfig()

	if err := database.Setup(conf); err != nil {
		logger.Fatalf("could not connect to database: %s", err)
	}

	if err := redis.Init(); err != nil {
		logger.Fatalf("could not connect to redis server: %s", err)
	}

	db := database.GetDB()
	rdb := redis.GetStore()

	r := router.Setup(db, rdb, conf)

	r.GET("/api/v1/docs/*any",
		ginSwagger.WrapHandler(swaggerFiles.Handler),
	)

	r.NoRoute(func(c *gin.Context) {
		c.JSON(http.StatusNotFound, gin.H{"code": "PAGE_NOT_FOUND", "message": "page not found"})
	})

	host := "0.0.0.0"
	port := "8080"

	if h := conf.Server.Host; h != "" {
		host = h
	}
	if p := conf.Server.Port; p > 0 {
		port = fmt.Sprint(p)
	}

	logger.Infof("Licensing Server starting at %s:%s", host, port)
	logger.Fatalf("%v", r.Run(host+":"+port))
}
