package httputil

import (
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"gorm.io/gorm"
)

type HTTPError struct {
	Code    int    `json:"code" example:"400"`
	Message string `json:"message" example:"status bad request"`
}

type ValidationError struct {
	FieldError validator.FieldError
}

func NewError(status int, err error) HTTPError {
	httpStatus := status
	if errors.Is(err, gorm.ErrRecordNotFound) {
		httpStatus = http.StatusNotFound
	}
	return HTTPError{
		Code:    httpStatus,
		Message: err.Error(),
	}
}

func (v ValidationError) NewFieldError() string {
	var sb strings.Builder

	sb.WriteString("validation failed on field '" + v.FieldError.Field() + "'")
	sb.WriteString(", condition: " + v.FieldError.ActualTag())

	if v.FieldError.Param() != "" {
		sb.WriteString(" { " + v.FieldError.Param() + " }")
	}

	if v.FieldError.Value() != nil && v.FieldError.Value() != "" {
		sb.WriteString(fmt.Sprintf(", actual: %v", v.FieldError.Value()))
	}

	return sb.String()

}

func HandleFieldError(ctx *gin.Context, err error) {
	for _, fieldErr := range err.(validator.ValidationErrors) {
		msg := ValidationError{FieldError: fieldErr}.NewFieldError()
		NewError(http.StatusBadRequest, errors.New(msg))
		return
	}
}
