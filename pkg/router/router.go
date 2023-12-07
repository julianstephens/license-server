package router

import (
	"net/http"

	"github.com/gin-contrib/gzip"
	"github.com/gin-gonic/gin"
	"github.com/julianstephens/license-server/controller"
	docs "github.com/julianstephens/license-server/docs"
	"github.com/julianstephens/license-server/pkg/logger"
	"github.com/julianstephens/license-server/pkg/middleware"
	sloggin "github.com/samber/slog-gin"
	"gorm.io/gorm"
)

var apiLogger = logger.GetLogger()

func Setup(db *gorm.DB) *gin.Engine {
	r := gin.New()

	routerLogger := logger.GetLogger()

	r.Use(sloggin.New(routerLogger))
	r.Use(gin.Recovery())

	r.Use(gzip.Gzip(gzip.DefaultCompression))

	docs.SwaggerInfo.BasePath = "/api/v1"

	api := controller.Controller{DB: db, Logger: apiLogger}

	publicGroup := r.Group("/api/v1")
	{
		publicGroup.GET("/", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{
				"name":    "Licensing Server API",
				"version": "0.1.0",
				"status":  "healthy",
			})
		})

		auth := publicGroup.Group("/auth")
		{
			auth.POST("/register", api.Register)
			auth.POST("/token", api.CreateToken)
		}
	}

	protectedGroup := r.Group("/api/v1", middleware.AuthGuard(api))
	{
		admin := protectedGroup.Group("/admin")
		{
			admin.POST("/users", api.AddUser)
		}
	}

	return r
}
