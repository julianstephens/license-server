package controller

import (
	"errors"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/julianstephens/license-server/model"
	"github.com/julianstephens/license-server/pkg/httputil"
	"github.com/julianstephens/license-server/service"
)

// GetUsers godoc
// @Summary Get all users
// @Description retrieves all users
// @Tags users
// @Security ApiKey
// @Success 200 {object} httputil.HTTPResponse[[]model.User]
// @Failure 500 {object} httputil.HTTPError
// @Router /admin/users [get]
func (base *Controller) GetUsers(c *gin.Context) {
	res, err := service.GetAll[model.User](base.DB)
	if err != nil {
		httputil.NewError(c, http.StatusInternalServerError, err)
		return
	}

	httputil.NewResponse(c, http.MethodGet, res)
}

// GetUser godoc
// @Summary Get a user
// @Description retrieve a specific user
// @Tags users
// @Security ApiKey
// @Success 200 {object} httputil.HTTPResponse[model.User]
// @Failure 400 {object} httputil.HTTPError
// @Failure 500 {object} httputil.HTTPError
// @Router /admin/users/:id [get]
func (base *Controller) GetUser(c *gin.Context) {
	userId := c.Param("id")
	if userId == "" {
		httputil.NewError(c, http.StatusBadRequest, errors.New("no user id provided"))
		return
	}

	res, err := service.FindById[model.User](base.DB, userId)
	if err != nil {
		httputil.NewError(c, http.StatusInternalServerError, err)
		return
	}

	httputil.NewResponse(c, http.MethodGet, res)
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
	var user model.User

	if err := c.ShouldBindJSON(&user); err != nil {
		httputil.NewError(c, http.StatusBadRequest, err)
		return
	}

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
		httputil.NewError(c, http.StatusBadRequest, err)
		return
	}

	if err := c.ShouldBindJSON(&user); err != nil {
		httputil.NewError(c, http.StatusBadRequest, err)
		return
	}

	res, err := service.Update[model.User](base.DB, userId, user)
	if err != nil {
		httputil.NewError(c, http.StatusInternalServerError, err)
		return
	}

	httputil.NewResponse(c, http.MethodPut, res)
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
		httputil.NewError(c, http.StatusBadRequest, err)
		return
	}

	res, err := service.Delete[model.User](base.DB, userId, model.User{})
	if err != nil {
		httputil.NewError(c, http.StatusInternalServerError, err)
		return
	}

	httputil.NewResponse(c, http.MethodDelete, res)
}
