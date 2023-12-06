package controller

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/julianstephens/license-server/model"
	"github.com/julianstephens/license-server/pkg/httputil"
)

// AddUser godoc
// @Summary Add a user
// @Description add new user
// @Tags users
// @Success 200 {object} model.User
// @Failure 400 {object} httputil.HTTPError
// @Router /users [get]
func (base *Controller) AddUser(c *gin.Context) {
	var newUser model.User
	if err := c.ShouldBindJSON(&newUser); err != nil {
		httputil.NewError(c, http.StatusBadRequest, err)
		return
	}
}
