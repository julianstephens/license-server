package controller

import (
	"errors"
	"net/http"
	"strings"

	mapset "github.com/deckarep/golang-set/v2"
	"github.com/gin-gonic/gin"
	"github.com/julianstephens/license-server/model"
	"github.com/julianstephens/license-server/pkg/httputil"
	"github.com/julianstephens/license-server/service"
)

type AuthRequest struct {
	Name     string `binding:"omitempty,alpha"`
	Email    string `binding:"required,email"`
	Password string
}

// Register godoc
// @Summary Register a user
// @Description registers new application user
// @Tags auth
// @Param data body AuthRequest true "new user info"
// @Success 200 {object} httputil.HTTPResponse[model.User]
// @Failure 400 {object} httputil.HTTPError
// @Failure 500 {object} httputil.HTTPError
// @Router /auth/register [post]
func (base *Controller) Register(c *gin.Context) {
	var req AuthRequest
	var user model.User

	if err := c.ShouldBindJSON(&req); err != nil {
		httputil.HandleFieldError(c, err)
		return
	}

	hashedPassword, err := service.HashData(req.Password)
	if err != nil {
		httputil.NewError(c, http.StatusBadRequest, errors.New("unable to parse password"))
		return
	}

	user.Name = req.Name
	user.Email = req.Email
	user.Password = hashedPassword

	res, err := service.Create[model.User](base.DB, user)
	if err != nil {
		if strings.Contains(err.Error(), "duplicate") {
			httputil.NewError(c, http.StatusBadRequest, err)
		} else {
			httputil.NewError(c, http.StatusInternalServerError, err)
		}
		return
	}

	httputil.NewResponse(c, http.MethodPost, res)
}

// CreateToken godoc
// @Summary Create a token
// @Description creates a new API key
// @Tags auth
// @Param data body AuthRequest true "returning user info"
// @Success 200 {object} httputil.HTTPResponse[model.DisplayAPIKey]
// @Failure 400 {object} httputil.HTTPError
// @Failure 500 {object} httputil.HTTPError
// @Router /auth/token [post]
func (base *Controller) CreateToken(c *gin.Context) {
	var req AuthRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		httputil.HandleFieldError(c, err)
		return
	}

	user, err := service.FindUserByEmail(base.DB, req.Email)
	if err != nil {
		httputil.NewError(c, http.StatusBadRequest, err)
		return
	}

	isValidUser := service.CompareWithHash(req.Password, user.Password)
	if !isValidUser {
		httputil.NewError(c, http.StatusBadRequest, errors.New("invalid password"))
		return
	}

	key, err := service.GenerateKey(base.DB, user.ID)
	if err != nil {
		httputil.NewError(c, http.StatusInternalServerError, err)
		return
	}

	httputil.NewResponse(c, http.MethodGet, key)
}

func (base *Controller) Authorize(key string, scopes ...string) (bool, error) {
	apiKey, err := service.FindByKey(base.DB, key)
	if err != nil {
		return false, err
	}

	if len(scopes) > 0 && !mapset.NewSet(strings.Split(apiKey.Scopes, ",")...).Contains(scopes...) {
		return false, errors.New("user does not have required scopes")
	}

	return true, nil
}
