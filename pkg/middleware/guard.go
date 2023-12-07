package middleware

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/julianstephens/license-server/controller"
	"github.com/julianstephens/license-server/pkg/httputil"
)

func AuthGuard(api controller.Controller, scopes ...string) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		apiKey := ctx.GetHeader("X-API-KEY")
		if apiKey == "" {
			httputil.NewError(ctx, http.StatusUnauthorized, errors.New("no API key provided"))
			ctx.Abort()
			return
		}

		isAuthed, err := api.Authorize(apiKey, scopes...)

		if err != nil || !isAuthed {
			httputil.NewError(ctx, http.StatusUnauthorized, errors.New("unauthorized"))
			ctx.Abort()
			return
		}

		ctx.Next()
	}
}
