package controller

import (
	"log/slog"

	"gorm.io/gorm"
)

type Controller struct {
	DB     *gorm.DB
	Logger *slog.Logger
}
