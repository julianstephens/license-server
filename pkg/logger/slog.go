package logger

import (
	"log/slog"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/lmittmann/tint"
)

var Logger *slog.Logger

func Setup() {
	w := os.Stdout
	l := slog.New(tint.NewHandler(w, nil))
	l = l.With("gin_mode", gin.EnvGinMode)
	Logger = l
}

func GetLogger() *slog.Logger {
	return Logger
}
