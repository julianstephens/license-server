package main

import (
	"time"

	ginzap "github.com/gin-contrib/zap"
	"github.com/gin-gonic/gin"
	"github.com/julianstephens/license-server/internal/routers"
	"go.uber.org/zap"
)

func main() {
	r := gin.New()

	logger, _ := zap.NewProduction()

	r.Use(ginzap.Ginzap(logger, time.RFC3339, true))

	publicGroup := r.Group("/")
	routers.NewPingRouter(publicGroup)
	r.Run(":8080")
}
