package backend

import (
	"context"
	"fmt"
	"strings"

	"github.com/jellydator/ttlcache/v3"
	jsoniter "github.com/json-iterator/go"
	"gorm.io/gorm"

	"github.com/julianstephens/license-server/backend/pkg/authserver"
	"github.com/julianstephens/license-server/backend/pkg/config"
	"github.com/julianstephens/license-server/backend/pkg/database"
	"github.com/julianstephens/license-server/backend/pkg/logger"
	"github.com/julianstephens/license-server/backend/pkg/model"
	"github.com/julianstephens/license-server/backend/pkg/service"
)

// App struct
type App struct {
	ctx   context.Context
	cache *ttlcache.Cache[string, string]
	DB    *gorm.DB
	Conf  *model.Config
}

// NewApp creates a new App application struct
func NewApp() *App {
	return &App{}
}

// startup is called at application startup
func (a *App) Startup(ctx context.Context) {
	// Perform your setup here
	logPath := "ls.log"
	logger.SetLogFile(&logPath)
	a.ctx = ctx
	a.cache = ttlcache.New[string, string]()
	a.DB = database.GetDB()
	a.Conf = config.GetConfig()
	logger.Infof("License Server starting...")
}

// domReady is called after front-end resources have been loaded
func (a App) DomReady(ctx context.Context) {
	// Add your action here
}

// beforeClose is called when the application is about to quit,
// either by clicking the window close button or calling runtime.Quit.
// Returning true will cause the application to continue, false will continue shutdown as normal.
func (a *App) BeforeClose(ctx context.Context) (prevent bool) {
	return false
}

// shutdown is called at application termination
func (a *App) Shutdown(ctx context.Context) {
	// Perform your teardown here
}

// Register creates a new application user
func (a *App) Register(req model.AuthRequest) (model.User, error) {
	hashedPassword, err := service.HashPassword(req.Password)
	if err != nil {
		return *new(model.User), service.HandleError(err, "unable to parse password", nil)
	}

	var user model.User
	user.Name = req.Name
	user.Email = req.Email
	user.Password = hashedPassword

	returned, err := service.Create[model.User](a.DB, user)
	if err != nil {
		if strings.Contains(err.Error(), "duplicate") {
			return *new(model.User), service.HandleError(err, "a user already exists with email address: %s", &[]any{user.Email})
		} else {
			return *new(model.User), service.HandleError(err, "unable to create user", nil)
		}
	}

	logger.Infof("registered user with email: %s", user.Email)
	return *returned, nil
}

func (a *App) SignIn() (bool, error) {
	isAuthed, err := authserver.Guard(a.Conf)
	if err != nil {
		return false, service.HandleError(err, "unable to authenticate", nil)
	}

	if isAuthed {
		logger.Infof("user is authenticated, attempting navigate")
		return true, nil
	}

	return false, nil
}

func (a *App) SignOut() (bool, error) {
	_, tokenPath, err := service.GetSecureFilePath(a.Conf, service.NewTrue())
	if err != nil {
		return false, err
	}

	if err := service.ClearToken(tokenPath); err != nil {
		return false, err
	}

	a.cache.Delete("profile")

	return true, nil
}

// CreateToken creates a new API key for a given user
func (a *App) CreateToken(req model.AuthRequest) (model.DisplayAPIKey, error) {
	user, err := service.FindUserByEmail(a.DB, req.Email)
	if err != nil {
		return *new(model.DisplayAPIKey), service.HandleError(err, "unable to find user with email: %s", &[]any{req.Email})
	}

	isValidUser := service.CompareWithHash(req.Password, user.Password)
	if !isValidUser {
		errMsg := "provided password is invalid"
		return *new(model.DisplayAPIKey), service.HandleError(fmt.Errorf(errMsg), errMsg, nil)
	}

	key, err := service.GenerateKey(a.DB, user.ID)
	if err != nil {
		return *new(model.DisplayAPIKey), service.HandleError(err, "unable to generate new authentication token", nil)
	}

	logger.Infof("created api key for user: %s", user.Email)
	return *key, nil
}

func (a *App) WhoAmI() (map[string]any, error) {
	_, tokenPath, err := service.GetSecureFilePath(a.Conf, service.NewTrue())
	if err != nil {
		return nil, err
	}

	token, _, err := service.GetAuthToken(a.Conf, tokenPath)
	if err != nil {
		a.cache.Delete("profile")
		return nil, err
	}

	if a.cache.Has("profile") {
		var profile map[string]any
		profileStr := a.cache.Get("profile")
		if err = jsoniter.UnmarshalFromString(profileStr.Value(), &profile); err == nil {
			return profile, nil
		}
	}

	profile, err := service.GetUserProfile(token, tokenPath, a.Conf, a.cache)
	if err == nil {
		profileStr, err := jsoniter.MarshalToString(profile)
		if err == nil {
			a.cache.Set("profile", profileStr, 0)
		}
	} else {
		a.cache.Delete("profile")
	}

	return profile, err
}
