package main

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	"github.com/julianstephens/license-server/internal/config"
	"github.com/julianstephens/license-server/internal/router"
	"github.com/julianstephens/license-server/pkg/httputil"
	"github.com/julianstephens/license-server/pkg/logger"
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
	logger.SetLogFile(nil)
	r := router.Setup()

	r.GET("/api/v1/docs/*any",
		ginSwagger.WrapHandler(swaggerFiles.Handler),
	)

	r.NoRoute(func(c *gin.Context) {
		httputil.NewError(c, http.StatusNotFound, fmt.Errorf("resource not found"))
	})

	conf := config.GetConfig()
	host := "0.0.0.0"
	port := "8080"

	if h := conf.Server.Host; h != "" {
		host = h
	}
	if p := conf.Server.Port; p > 0 {
		port = strconv.Itoa(p)
	}

	logger.Infof("Licensing Server starting at %s:%s", host, port)
	logger.Fatalf("%v", r.Run(host+":"+port))
}
