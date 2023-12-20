package controller

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/julianstephens/license-server/backend/pkg/config"
	"github.com/julianstephens/license-server/backend/pkg/database"
	"github.com/julianstephens/license-server/backend/pkg/httputil"
	"github.com/julianstephens/license-server/backend/pkg/licensemanager"
	"github.com/julianstephens/license-server/backend/pkg/model"
)

var (
	conf = config.GetConfig()
	db   = database.GetDB()
	lm   = licensemanager.LicenseManager{Config: conf, DB: db}
)

// IssueLicense godoc
// @Summary Issue a license
// @Description issues a new product license
// @Tags licenses
// @Security ApiKey
// @Param data body model.LicenseRequest true "new license info"
// @Success 200 {object} httputil.HTTPResponse[model.ActivationData]
// @Failure 400 {object} httputil.HTTPError
// @Failure 500 {object} httputil.HTTPError
// @Router /licenses/issue [post]
func (base *Controller) IssueLicense(c *gin.Context) {
	var req model.LicenseRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		httputil.NewError(http.StatusBadRequest, err)
		return
	}

	ok, license, err := lm.ValidateKey(req.Key)
	if err != nil || !ok {
		httputil.NewResponse[any](c, gin.H{"success": false, "data": nil}, httputil.Options{IsCrudHandler: false})
		return
	}

	res, err := lm.AssignLicense(req.Machine, license)
	if err != nil {
		if err != nil {
			httputil.NewError(http.StatusInternalServerError, err)
			return
		}
	}

	httputil.NewResponse[any](c, gin.H{"success": true, "data": &res}, httputil.Options{IsCrudHandler: false})
}

// RevokeLicense godoc
// @Summary Revoke a license
// @Description revokes a license with id
// @Tags licenses
// @Security ApiKey
// @Param data body model.LicenseRevokeRequest true "license id"
// @Success 204
// @Failure 400 {object} httputil.HTTPError
// @Failure 500 {object} httputil.HTTPError
// @Router /licenses/:id/revoke [delete]
func (base *Controller) RevokeLicense(c *gin.Context) {
	var req model.LicenseRevokeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		httputil.NewError(http.StatusBadRequest, err)
		return
	}

	err := lm.RevokeLicense(req.Id)
	if err != nil {
		httputil.NewError(http.StatusInternalServerError, err)
		return
	}

	httputil.NewResponse[any](c, nil, httputil.Options{IsCrudHandler: false, Status: http.StatusNoContent})
}
