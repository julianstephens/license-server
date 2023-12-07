package httputil

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

type HTTPResponse[T any] struct {
	Message string `json:"message" `
	Data    T      `json:"data"`
}

func NewResponse[T any](ctx *gin.Context, method string, data T) {
	status := http.StatusOK
	var message string

	switch strings.ToUpper(method) {
	case http.MethodPost:
		status = http.StatusCreated
		message = "resource(s) created successfully"
	case http.MethodPut:
		message = "resource(s) updated successfully"
	case http.MethodDelete:
		message = "resource(s) deleted successfully"
	default:
		message = "resource(s) retrieved successfully"
	}

	res := HTTPResponse[T]{
		Message: message,
		Data:    data,
	}
	ctx.JSON(status, res)
}
