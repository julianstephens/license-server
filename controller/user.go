package controller

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/julianstephens/license-server/model"
	"github.com/julianstephens/license-server/pkg/httputil"
	"github.com/julianstephens/license-server/service"
)

// AddUser godoc
// @Summary Add a user
// @Description add new user
// @Tags users
// @Security ApiKey
// @Param data body model.User true "new user info"
// @Success 200 {object} model.User
// @Failure 400 {object} httputil.HTTPError
// @Failure 500 {object} httputil.HTTPError
// @Router /admin/users [post]
func (base *Controller) AddUser(c *gin.Context) {
	var user *model.User

	if err := c.ShouldBindJSON(&user); err != nil {
		httputil.NewError(c, http.StatusBadRequest, err)
		return
	}

	res, err := service.CreateUser(base.DB, user)
	if err != nil {
		httputil.NewError(c, http.StatusInternalServerError, err)
	}

	httputil.NewResponse(c, http.MethodGet, res)
}
