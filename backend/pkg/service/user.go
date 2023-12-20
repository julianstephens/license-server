package service

import (
	"github.com/julianstephens/license-server/backend/pkg/model"
	"gorm.io/gorm"
)

func FindUserByEmail(db *gorm.DB, email string) (*model.User, error) {
	var result model.User
	if err := db.Where(&model.User{Email: email}).First(&result).Error; err != nil {
		return nil, err
	}

	return &result, nil
}
