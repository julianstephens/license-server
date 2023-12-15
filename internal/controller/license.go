package controller

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/julianstephens/license-server/internal/config"
	"github.com/julianstephens/license-server/internal/licensemanager"
	"github.com/julianstephens/license-server/internal/service"
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
	id, err := service.GetId(c)
	if err != nil {
		httputil.NewError(c, http.StatusInternalServerError, err)
		return
	}
	_, _, key, err := lm.GenerateLicense(id)
	if err != nil {
		httputil.NewError(c, http.StatusInternalServerError, err)
		return
	}
	httputil.NewResponse[any](c, http.MethodGet, gin.H{"product_key": key, "len": len(key)}, &httputil.Opts{IsCrudHandler: false})
}

type Req struct {
	Key     string
	Machine *string
}

// ValidateLicense godoc
// @Summary Validate a license
// @Description validates a product license
// @Tags licenses
// @Security ApiKey
// @Success 200 {object} httputil.HTTPResponse[any]
// @Failure 400 {object} httputil.HTTPError
// @Failure 500 {object} httputil.HTTPError
// @Router /licenses/validate [post]
func (base *Controller) ValidateLicense(c *gin.Context) {
	lm.SetCurrentUser(c.GetString("user"))
	var data Req
	if err := c.ShouldBindJSON(&data); err != nil {
		httputil.NewError(c, http.StatusInternalServerError, err)
		return
	}
	ok, err := lm.ValidateKey(data.Key)
	if err != nil {
		httputil.NewResponse[any](c, http.MethodPost, gin.H{"is_valid": false}, &httputil.Opts{IsCrudHandler: false})
		return
	}
	httputil.NewResponse[any](c, http.MethodPost, gin.H{"is_valid": ok}, &httputil.Opts{IsCrudHandler: false})
}
