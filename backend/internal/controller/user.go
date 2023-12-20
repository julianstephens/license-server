package controller

import (
	"errors"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	jsoniter "github.com/json-iterator/go"

	"github.com/julianstephens/license-server/backend/pkg/httputil"
	"github.com/julianstephens/license-server/backend/pkg/logger"
	"github.com/julianstephens/license-server/backend/pkg/model"
	"github.com/julianstephens/license-server/backend/pkg/service"
)

type ScopeRequest struct {
	Add    []string
	Remove []string
}

// GetUsers godoc
// @Summary Get all users
// @Description retrieves all users
// @Tags users
// @Security ApiKey
// @Success 200 {object} httputil.HTTPResponse[[]model.User]
// @Failure 500 {object} httputil.HTTPError
// @Router /admin/users [get]
func (base *Controller) GetUsers(c *gin.Context) {
	logger.Infof("retrieving all users as user <%s>", c.GetString("user"))
	res, err := service.GetAll[model.User](base.DB)
	if err != nil {
		httputil.NewError(http.StatusInternalServerError, err)
		return
	}

	httputil.NewResponse(c, res, httputil.Options{IsCrudHandler: true, HttpMsgMethod: httputil.Get})
}

// GetUser godoc
// @Summary Get a user
// @Description retrieve a specific user
// @Tags users
// @Security ApiKey
// @Success 200 {object} httputil.HTTPResponse[model.UserWithScopes]
// @Failure 400 {object} httputil.HTTPError
// @Failure 500 {object} httputil.HTTPError
// @Router /admin/users/:id [get]
func (base *Controller) GetUser(c *gin.Context) {
	var res model.UserWithScopes

	userId := c.Param("id")
	if userId == "" {
		httputil.NewError(http.StatusBadRequest, errors.New("no user id provided"))
		return
	}
	logger.Infof("retrieving user <%s> as user <%s>", userId, c.GetString("user"))

	user, err := service.FindById[model.User](base.DB, userId)
	if err != nil {
		httputil.NewError(http.StatusInternalServerError, err)
		return
	}
	res.User = *user

	key, _ := service.Find[model.APIKey](base.DB, model.APIKey{UserId: user.ID}, nil)
	scopes := key.Scopes
	if len(scopes) > 0 {
		res.Scopes = strings.Split(scopes, ",")
	}

	httputil.NewResponse(c, res, httputil.Options{IsCrudHandler: true, HttpMsgMethod: httputil.Get})
}

// AddUser godoc
// @Summary Add a user
// @Description creates a new user
// @Tags users
// @Security ApiKey
// @Param data body model.User true "new user info"
// @Success 200 {object} httputil.HTTPResponse[model.User]
// @Failure 400 {object} httputil.HTTPError
// @Failure 500 {object} httputil.HTTPError
// @Router /admin/users [post]
func (base *Controller) AddUser(c *gin.Context) {
	logger.Infof("creating new user as user <%s>", c.GetString("user"))
	var user model.User

	if err := c.ShouldBindJSON(&user); err != nil {
		httputil.NewError(http.StatusBadRequest, err)
		return
	}

	res, err := service.Create[model.User](base.DB, user)
	if err != nil {
		if strings.Contains(err.Error(), "duplicate") {
			httputil.NewError(http.StatusBadRequest, err)
		} else {
			httputil.NewError(http.StatusInternalServerError, err)
		}
		return
	}

	httputil.NewResponse(c, res, httputil.Options{IsCrudHandler: true, HttpMsgMethod: httputil.Post})
}

// UpdateUser godoc
// @Summary Update a user
// @Description updates a specific user
// @Tags users
// @Security ApiKey
// @Param data body model.User true "updated user info"
// @Success 200 {object} httputil.HTTPResponse[model.User]
// @Failure 400 {object} httputil.HTTPError
// @Failure 500 {object} httputil.HTTPError
// @Router /admin/users/:id [put]
func (base *Controller) UpdateUser(c *gin.Context) {
	var user model.User

	userId, err := service.GetId(c)
	if err != nil {
		httputil.NewError(http.StatusBadRequest, err)
		return
	}
	logger.Infof("updating user <%s> as user <%s>", userId, c.GetString("user"))

	if err := c.ShouldBindJSON(&user); err != nil {
		httputil.NewError(http.StatusBadRequest, err)
		return
	}

	var decoded map[string]interface{}
	b, err := jsoniter.Marshal(user)
	if err != nil {
		httputil.NewError(http.StatusInternalServerError, err)
		return
	}
	err = jsoniter.Unmarshal(b, &decoded)
	if err != nil {
		httputil.NewError(http.StatusInternalServerError, err)
		return
	}

	res, err := service.Update[model.User](base.DB, userId, decoded)
	if err != nil {
		httputil.NewError(http.StatusInternalServerError, err)
		return
	}

	httputil.NewResponse(c, res, httputil.Options{IsCrudHandler: true, HttpMsgMethod: httputil.Put})
}

// DeleteUser godoc
// @Summary Delete a user
// @Description deletes a specific user
// @Tags users
// @Security ApiKey
// @Success 200 {object} httputil.HTTPResponse[model.User]
// @Failure 400 {object} httputil.HTTPError
// @Failure 500 {object} httputil.HTTPError
// @Router /admin/users/:id [delete]
func (base *Controller) DeleteUser(c *gin.Context) {
	userId, err := service.GetId(c)
	if err != nil {
		httputil.NewError(http.StatusBadRequest, err)
		return
	}
	logger.Infof("deleting user <%s> as user <%s>", userId, c.GetString("user"))

	res, err := service.Delete[model.User](base.DB, userId, model.User{})
	if err != nil {
		httputil.NewError(http.StatusInternalServerError, err)
		return
	}

	httputil.NewResponse(c, res, httputil.Options{IsCrudHandler: true, HttpMsgMethod: httputil.Delete})
}

// UpdateScope godoc
// @Summary Update a user's scope
// @Description update scopes for a specific user
// @Tags users
// @Param data body ScopeRequest true "scopes to modify"
// @Success 200 {object} httputil.HTTPResponse[model.UserWithScopes]
// @Failure 400 {object} httputil.HTTPError
// @Failure 500 {object} httputil.HTTPError
// @Router /admin/users/:id/scopes [put]
func (base *Controller) UpdateUserScopes(c *gin.Context) {
	var req ScopeRequest

	userId, err := service.GetId(c)
	if err != nil {
		httputil.NewError(http.StatusBadRequest, err)
		return
	}
	logger.Infof("updating scopes for user <%s> as user <%s>", userId, c.GetString("user"))

	if err := c.ShouldBindJSON(&req); err != nil {
		httputil.HandleFieldError(c, err)
		return
	}

	key, err := service.Find[model.APIKey](base.DB, model.APIKey{UserId: userId}, &service.PreloadOptions{ShouldPreload: true, PreloadQuery: "User"})
	if err != nil {
		httputil.NewError(http.StatusBadRequest, err)
		return
	}

	user := key.User

	toAdd := service.If(len(req.Add) > 0, req.Add, []string{})
	toRemove := service.If(len(req.Remove) > 0, req.Remove, []string{})

	curScopes := strings.Split(key.Scopes, ",")
	cleanedScopes := service.Difference(curScopes, toRemove)
	updatedScopes := append(cleanedScopes, toAdd...)
	key.Scopes = strings.Join(updatedScopes, ",")

	res, err := service.Update[model.APIKey](base.DB, key.ID, map[string]any{"Scopes": strings.Join(updatedScopes, ",")})
	if err != nil {
		httputil.NewError(http.StatusInternalServerError, err)
		return
	}

	httputil.NewResponse(c, model.UserWithScopes{User: user, Scopes: strings.Split(res.Scopes, ",")}, httputil.Options{IsCrudHandler: true, HttpMsgMethod: httputil.Put})
}
