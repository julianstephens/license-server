package service

import (
	"errors"
	"strings"

	"github.com/google/uuid"
	"github.com/julianstephens/license-server/model"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

const KeyByteSize int = 32
const BcryptCost = 14

func HashData(data any) (string, error) {
	var bytes []byte
	var err error

	switch v := data.(type) {
	case string:
		bytes, err = bcrypt.GenerateFromPassword([]byte(v), BcryptCost)
	case []byte:
		bytes, err = bcrypt.GenerateFromPassword(v, BcryptCost)
	default:
		return "", errors.New("invalid data type")
	}

	return string(bytes), err
}

func CompareWithHash(data string, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(data))
	return err == nil
}

func GenerateKey(db *gorm.DB, userID string) (*model.APIKey, error) {
	if userID == "" {
		return nil, errors.New("no user provided")
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
	return &model.APIKey{UserId: userID, Key: apiKey, Mask: apiKey[:6]}, err
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
