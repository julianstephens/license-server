package router

import (
	"net/http"

	"github.com/gin-contrib/gzip"
	"github.com/gin-gonic/gin"
	"github.com/julianstephens/license-server/pkg/logger"
	sloggin "github.com/samber/slog-gin"
	"gorm.io/gorm"
)

func Setup(db *gorm.DB) *gin.Engine {
	r := gin.New()

	routerLogger := logger.GetLogger()

	r.Use(sloggin.New(routerLogger))
	r.Use(gin.Recovery())

	r.Use(gzip.Gzip(gzip.DefaultCompression))

	// api := controller.Controller{DB: db}

	publicGroup := r.Group("/api/v1")
	{
		publicGroup.GET("/", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{
				"name":    "Licensing Server API",
				"version": "0.1.0",
				"status":  "healthy",
			})
		})
	}

	return r
}
