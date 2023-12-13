package controller

import (
	"github.com/julianstephens/license-server/internal/model"
	"gorm.io/gorm"
)

type Controller struct {
	DB     *gorm.DB
	Config *model.Config
}
