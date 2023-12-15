package controller

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/julianstephens/license-server/internal/config"
	"github.com/julianstephens/license-server/internal/licensemanager"
	"github.com/julianstephens/license-server/internal/model"
	"github.com/julianstephens/license-server/pkg/database"
	"github.com/julianstephens/license-server/pkg/httputil"
)

var conf = config.GetConfig()
var db = database.GetDB()
var lm = licensemanager.LicenseManager{Config: conf, DB: db}

// IssueLicense godoc
// @Summary Issue a license
// @Description issues a new product license
// @Tags licenses
// @Security ApiKey
// @Param data body model.LicenseRequest true "new license info"
// @Success 200 {object} httputil.HTTPResponse[any]
// @Failure 400 {object} httputil.HTTPError
// @Failure 500 {object} httputil.HTTPError
// @Router /licenses/issue [post]
func (base *Controller) IssueLicense(c *gin.Context) {
	lm.SetCurrentUser(c.GetString("user"))
	var req model.LicenseRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		httputil.NewError(c, http.StatusBadRequest, err)
		return
	}

	ok, err := lm.ValidateKey(req.Key)
	if err != nil {
		httputil.NewResponse[any](c, gin.H{"success": false}, httputil.Options{IsCrudHandler: false})
		return
	}
	httputil.NewResponse[any](c, gin.H{"success": ok}, httputil.Options{IsCrudHandler: false})
}
