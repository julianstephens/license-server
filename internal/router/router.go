package router

import (
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/gin-contrib/gzip"
	"github.com/gin-gonic/gin"

	docs "github.com/julianstephens/license-server/docs"
	"github.com/julianstephens/license-server/internal/config"
	"github.com/julianstephens/license-server/internal/controller"
	"github.com/julianstephens/license-server/internal/middleware"
	"github.com/julianstephens/license-server/pkg/database"
	"github.com/julianstephens/license-server/pkg/logger"
)

const BasePath = "/api/v1"

func Setup() *gin.Engine {
	r := gin.New()

	f, err := os.OpenFile("ls.access.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		logger.Errorf("failed to create access log file: %v", err)
	} else {
		gin.DefaultWriter = io.MultiWriter(f)
	}

	r.Use(gin.Logger())
	r.Use(gin.Recovery())
	r.Use(gzip.Gzip(gzip.DefaultCompression))

	docs.SwaggerInfo.BasePath = BasePath

	db := database.GetDB()
	conf := config.GetConfig()

	api := controller.Controller{DB: db, Config: conf}

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
			licenseGroup.GET("/issue/:id", api.IssueLicense)
			licenseGroup.POST("/validate/:id", api.ValidateLicense)
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
			productGroup.GET("/:id/key", api.CreateProductKeyPair)
			productGroup.POST("/", api.AddProduct)
			productGroup.PUT("/:id", api.UpdateProduct)
			productGroup.DELETE("/:id", api.DeleteProduct)
		}
	}

	return r
}
