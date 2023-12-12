package router

import (
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/chenyahui/gin-cache/persist"
	"github.com/gin-contrib/gzip"
	"github.com/gin-gonic/gin"
	"github.com/julianstephens/license-server/controller"
	docs "github.com/julianstephens/license-server/docs"
	"github.com/julianstephens/license-server/pkg/logger"
	"github.com/julianstephens/license-server/pkg/middleware"
	sloggin "github.com/samber/slog-gin"
	"gorm.io/gorm"
)

const BasePath = "/api/v1"

var apiLogger = logger.GetLogger()

func Setup(db *gorm.DB, rdb *persist.RedisStore) *gin.Engine {
	r := gin.New()

	f, err := os.OpenFile("ls.access.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		logger.Errorf("Failed to create access log file: %v", err)
	} else {
		gin.DefaultWriter = io.MultiWriter(f)
	}

	routerLogger := logger.GetLogger()

	r.Use(sloggin.New(routerLogger))
	r.Use(gin.Recovery())

	r.Use(gzip.Gzip(gzip.DefaultCompression))

	docs.SwaggerInfo.BasePath = BasePath

	api := controller.Controller{DB: db, Logger: apiLogger}

	publicGroup := r.Group(BasePath)
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

	protectedGroup := r.Group(BasePath, middleware.AuthGuard(api))
	{
		licenseGroup := protectedGroup.Group("/licenses")
		{
			licenseGroup.GET("/issue", api.IssueLicense)
			licenseGroup.POST("/validate", api.ValidateLicense)
		}
	}

	adminGroup := r.Group(fmt.Sprintf("%s/admin", BasePath), middleware.AuthGuard(api, "admin"))
	{
		userGroup := adminGroup.Group("/users")
		{
			userGroup.GET("/", api.GetUsers)
			userGroup.GET("/:id", api.GetUser)
			userGroup.POST("/", api.AddUser)
			userGroup.PUT("/:id", api.UpdateUser)
			userGroup.DELETE("/:id", api.DeleteUser)
			userGroup.PUT("/:id/scopes", api.UpdateUserScopes)
		}

		productGroup := adminGroup.Group("/products")
		{
			productGroup.GET("/", api.GetProducts)
			productGroup.GET("/:id", api.GetProduct)
			productGroup.POST("/", api.AddProduct)
			productGroup.PUT("/:id", api.UpdateProduct)
			productGroup.DELETE("/:id", api.DeleteProduct)
		}
	}

	return r
}
