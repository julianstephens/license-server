package controller

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/julianstephens/license-server/internal/controller/httputil"
)

// IssueLicense godoc
// @Summary Issue a license
// @Description issues a new product license
// @Tags licenses
// @Security ApiKey
// @Success 200 {object} httputil.HTTPResponse[any]
// @Failure 400 {object} httputil.HTTPError
// @Failure 500 {object} httputil.HTTPError
// @Router /licenses/issue [get]
func (base *Controller) IssueLicense(c *gin.Context) {
	httputil.NewResponse[any](c, http.MethodGet, gin.H{"message": "ok"})
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
	httputil.NewResponse[any](c, http.MethodPost, gin.H{"message": "ok"})
}
