package routers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func ping(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "pong"})
}

func NewPingRouter(group *gin.RouterGroup) {
	group.GET("/ping", ping)
}
