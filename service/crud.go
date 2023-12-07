package service

import (
	"gorm.io/gorm"
)

func GetAll[T any](db *gorm.DB) (*[]T, error) {
	var result []T

	if err := db.Find(&result).Error; err != nil {
		return nil, err
	}

	return &result, nil
}

func FindById[T any](db *gorm.DB, id string) (*T, error) {
	var result T
	if err := db.Where("id = ?", id).First(&result).Error; err != nil {
		return nil, err
	}

	return &result, nil
}

func Create[T any](db *gorm.DB, newResource T, conditions ...T) (*T, error) {
	if err := db.Create(&newResource).Error; err != nil {
		return nil, err
	}

	return &newResource, nil
}

func Update[T any](db *gorm.DB, id string, updates T) (*T, error) {
	existing, err := FindById[T](db, id)
	if err != nil {
		return nil, err
	}

	// updates non-null/empty fields in existing resource
	items := StructItems(updates)
	for _, item := range items {
		if item.Value != "" && item.Value != nil {
			SetProperty(existing, item.Key, item.Value)
		}
	}

	if err := db.Save(&existing).Error; err != nil {
		return nil, err
	}

	return existing, nil
}

func Delete[T any](db *gorm.DB, id string, res T) (*T, error) {
	var result T
	if err := db.Delete(res, id).Error; err != nil {
		return nil, err
	}

	return &result, nil
}
