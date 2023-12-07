package service

import (
	"github.com/julianstephens/license-server/model"
	"github.com/julianstephens/license-server/pkg/logger"
	"gorm.io/gorm"
)

var userLogger = logger.GetLogger()

func CreateUser(db *gorm.DB, user *model.User) (*model.User, error) {
	res := db.FirstOrCreate(&user, model.User{Email: user.Email})
	if res.Error != nil {
		userLogger.Error("CreateUser error: %v", res.Error)
		return user, res.Error
	}

	if res.RowsAffected > 0 {
		userLogger.Error("CreateUser error: user already exists")
		return user, res.Error
	}

	return user, nil
}

func FindByEmail(db *gorm.DB, email string) (*model.User, error) {
	var result model.User
	if err := db.Where(&model.User{Email: email}).First(&result).Error; err != nil {
		return nil, err
	}

	return &result, nil
}
