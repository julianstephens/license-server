package service

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/julianstephens/license-server/pkg/model"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

const KeyByteSize int = 32
const BcryptCost int = 14

func HashData(data any) ([]byte, error) {
	var bytes []byte
	var err error

	switch v := data.(type) {
	case string:
		bytes, err = bcrypt.GenerateFromPassword([]byte(v), BcryptCost)
	case []byte:
		bytes, err = bcrypt.GenerateFromPassword(v, BcryptCost)
	default:
		return []byte{}, errors.New("invalid data type")
	}

	return bytes, err
}

func CompareWithHash(data string, hash []byte) bool {
	err := bcrypt.CompareHashAndPassword(hash, []byte(data))
	return err == nil
}

func GenerateKey(db *gorm.DB, userID string) (*model.DisplayAPIKey, error) {
	if userID == "" {
		return nil, errors.New("no user provided")
	}

	existingKey, _ := Find[model.APIKey](db, model.APIKey{UserId: userID}, nil)
	if existingKey != nil {
		return nil, fmt.Errorf("found existing api key: %s", existingKey.Mask)
	}

	uid := uuid.New().String()
	apiKey := strings.Join(strings.Split(uid, "-"), "")

	hashedKey, err := HashData(apiKey)
	if err != nil {
		return nil, err
	}

	key := &model.APIKey{
		UserId: userID,
		Key:    hashedKey,
		Mask:   apiKey[:6],
	}

	err = db.Save(key).Error
	return &model.DisplayAPIKey{Base: key.Base, UserId: userID, Key: apiKey, Mask: apiKey[:6], ExpiresAt: time.Unix(key.ExpiresAt, 0).Local().String()}, err
}

func FindByKey(db *gorm.DB, key string) (*model.APIKey, error) {
	var apiKey model.APIKey

	err := db.Where(&model.APIKey{Mask: key[:6]}).First(&apiKey).Error
	if err != nil {
		return nil, err
	}

	isValidKey := CompareWithHash(key, apiKey.Key)
	if !isValidKey {
		return nil, errors.New("invalid key")
	}

	return &apiKey, nil
}

func FindByOwner(db *gorm.DB, key string) (*model.APIKey, error) {
	var apiKey model.APIKey

	err := db.Where(&model.APIKey{Mask: key[:6]}).First(&apiKey).Error
	if err != nil {
		return nil, err
	}

	isValidKey := CompareWithHash(key, apiKey.Key)
	if !isValidKey {
		return nil, errors.New("invalid key")
	}

	return &apiKey, nil
}
