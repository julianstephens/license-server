package controller

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/julianstephens/license-server/model"
	"github.com/julianstephens/license-server/pkg/httputil"
	"github.com/julianstephens/license-server/service"
)

type LoginRequest struct {
	Email    string
	Password string
}

// Register godoc
// @Summary Registers a user
// @Description register new application user
// @Tags auth
// @Param data body model.User true "new user info"
// @Success 200 {object} httputil.HTTPResponse[model.User]
// @Failure 400 {object} httputil.HTTPError
// @Failure 500 {object} httputil.HTTPError
// @Router /auth/register [post]
func (base *Controller) Register(c *gin.Context) {
	user := new(model.User)
	if err := c.ShouldBindJSON(&user); err != nil {
		for _, fieldErr := range err.(validator.ValidationErrors) {
			msg := httputil.ValidationError{FieldError: fieldErr}.NewFieldError()
			httputil.NewError(c, http.StatusBadRequest, errors.New(msg))
			return
		}
	}

	hashedPassword, err := service.HashData(user.Password)
	if err != nil {
		httputil.NewError(c, http.StatusBadRequest, errors.New("unable to parse password"))
		return
	}

	user.Password = hashedPassword

	res, err := service.CreateUser(base.DB, user)
	if err != nil {
		httputil.NewError(c, http.StatusInternalServerError, err)
	}

	httputil.NewResponse(c, http.MethodPost, res)
}

// CreateToken godoc
// @Summary Creates a token
// @Description creates a new API key
// @Tags auth
// @Param data body LoginRequest true "returning user info"
// @Success 200 {object} httputil.HTTPResponse[model.APIKey]
// @Failure 400 {object} httputil.HTTPError
// @Failure 500 {object} httputil.HTTPError
// @Router /auth/token [post]
func (base *Controller) CreateToken(c *gin.Context) {
	req := new(LoginRequest)
	if err := c.ShouldBindJSON(&req); err != nil {
		httputil.NewError(c, http.StatusBadRequest, err)
		return
	}

	user, err := service.FindByEmail(base.DB, req.Email)
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

func (base *Controller) Authorize(key string) (bool, error) {
	apiKey, err := service.FindByKey(base.DB, key)
	return apiKey != nil, err
}
