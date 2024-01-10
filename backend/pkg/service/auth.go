package service

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/jellydator/ttlcache/v3"
	jsoniter "github.com/json-iterator/go"
	"github.com/labstack/gommon/log"
	"github.com/lestrrat-go/jwx/v2/jwk"
	"github.com/lestrrat-go/jwx/v2/jws"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"

	"github.com/julianstephens/license-server/backend/pkg/model"
)

const (
	KeyByteSize int = 32
	BcryptCost  int = 14
)

func HashPassword(data any) ([]byte, error) {
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

	hashedKey, err := HashPassword(apiKey)
	if err != nil {
		return nil, err
	}

	key := &model.APIKey{
		UserId: userID,
		Key:    hashedKey,
		Mask:   apiKey[:6],
	}

	err = db.Save(key).Error
	return &model.DisplayAPIKey{Base: key.Base, UserId: userID, Key: apiKey, Mask: apiKey[:6], ExpiresAt: time.Unix(key.ExpiresAt, 0).Local().Format(time.DateTime)}, err
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

func NewJWKSet(jwkUrl string) jwk.Set {
	jwkCache := jwk.NewCache(context.Background())

	// register a minimum refresh interval for this URL.
	// when not specified, defaults to Cache-Control and similar resp headers
	err := jwkCache.Register(jwkUrl, jwk.WithMinRefreshInterval(10*time.Minute))
	if err != nil {
		panic("failed to register jwk location")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// fetch once on application startup
	_, err = jwkCache.Refresh(ctx, jwkUrl)
	if err != nil {
		log.Fatalf("failed to refresh auth0 JWKS: %+v", err)
	}
	// create the cached key set
	return jwk.NewCachedSet(jwkCache, jwkUrl)
}

func VerifyToken(token []byte, signingKeys jwk.Set) (ok bool) {
	var err error
	_, err = jws.Verify(token, jws.WithKeySet(signingKeys, jws.WithInferAlgorithmFromKey(true)))
	if err == nil {
		ok = true
	} else {
		ok = false
	}

	return
}

func GetAuthToken(conf *model.Config, tokenPath string) (token []byte, signingKeys jwk.Set, err error) {
	signingKeys = NewJWKSet(conf.Auth.JwksUrl)

	tokenStr, err := LoadToken(tokenPath)
	if err != nil {
		return *new([]byte), signingKeys, err
	}

	ok := VerifyToken([]byte(tokenStr), signingKeys)
	if !ok {
		return *new([]byte), signingKeys, HandleError(fmt.Errorf("could not to verify stored token"), "unauthenticated", nil)
	}

	return []byte(tokenStr), signingKeys, nil
}

func GetUserProfile(token []byte, tokenPath string, conf *model.Config, cache *ttlcache.Cache[string, string]) (map[string]any, error) {
	req, err := http.NewRequest("GET", fmt.Sprintf("https://%s.us.auth0.com/userinfo", conf.Auth.TenantID), nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode == http.StatusUnauthorized {
		cache.Delete("profile")
		err = ClearToken(tokenPath)
		if err != nil {
			return nil, err
		}

		return nil, fmt.Errorf("unauthenticated")
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var json map[string]any
	err = jsoniter.Unmarshal(body, &json)
	if err != nil {
		return nil, err
	}

	return json, nil
}
