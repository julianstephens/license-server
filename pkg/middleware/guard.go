package middleware

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/julianstephens/license-server/controller"
	"github.com/julianstephens/license-server/pkg/httputil"
	"gorm.io/gorm"
)

func AuthGuard(db *gorm.DB) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		apiKey := ctx.GetHeader("X-API-KEY")
		if apiKey == "" {
			httputil.NewError(ctx, http.StatusUnauthorized, errors.New("no API key provided"))
			ctx.Abort()
			return
		}

		api := controller.Controller{DB: db}
		isAuthed, err := api.Authorize(apiKey)

		if err != nil || !isAuthed {
			httputil.NewError(ctx, http.StatusUnauthorized, errors.New("unauthorized"))
			ctx.Abort()
			return
		}

		ctx.Next()
	}
}
