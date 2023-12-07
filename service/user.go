package service

import (
	"github.com/julianstephens/license-server/model"
	"github.com/julianstephens/license-server/pkg/logger"
	"gorm.io/gorm"
)

var userLogger = logger.GetLogger()

func CreateUser(db *gorm.DB, user *model.User) (*model.User, error) {
	res := db.Where(model.User{Email: user.Email}).FirstOrCreate(&user)
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

func ModifyUser(db *gorm.DB, updatedUser *model.User, userId string) (*model.User, error) {
	existingUser, err := FindUserById(db, userId)
	if err != nil {
		return nil, err
	}

	// updates non-null/empty fields in existing user
	items := StructItems(updatedUser)
	for _, item := range items {
		if item.Value != "" && item.Value != nil {
			SetProperty(existingUser, item.Key, item.Value)
		}
	}

	if err := db.Save(&existingUser).Error; err != nil {
		return nil, err
	}

	return existingUser, nil
}

func DeleteUser(db *gorm.DB, userId string) (*model.User, error) {
	existingUser, err := FindUserById(db, userId)
	if err != nil {
		return nil, err
	}

	if err := db.Delete(&existingUser).Error; err != nil {
		return nil, err
	}

	return existingUser, nil
}

func FindUserByEmail(db *gorm.DB, email string) (*model.User, error) {
	var result model.User
	if err := db.Where(&model.User{Email: email}).First(&result).Error; err != nil {
		return nil, err
	}

	return &result, nil
}

func FindUserById(db *gorm.DB, userId string) (*model.User, error) {
	var result model.User
	if err := db.Where("id = ?", userId).First(&result).Error; err != nil {
		return nil, err
	}

	return &result, nil
}

func GetAllUsers(db *gorm.DB) (*[]model.User, error) {
	var result []model.User

	if err := db.Find(&result).Error; err != nil {
		return nil, err
	}

	return &result, nil
}
