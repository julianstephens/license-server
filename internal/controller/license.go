package controller

import (
	"crypto/elliptic"
	"crypto/x509"
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/julianstephens/license-server/internal/service"
	"github.com/julianstephens/license-server/pkg/crypto"
	"github.com/julianstephens/license-server/pkg/httputil"
	"github.com/julianstephens/license-server/pkg/logger"
)

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
	var ec = crypto.New(elliptic.P256())
	id, err := service.GetId(c)
	if err != nil {
		panic(err)
	}

	kp, err := service.LoadKey(id, base.Config)
	if err != nil {
		logger.Infof("unable to load key")
		panic(err)
	}
	logger.Infof("kp: %s", kp.PrivateKey)

	privKey, err := ec.DecodePrivate(kp.PrivateKey)
	if err != nil {
		logger.Infof("unable to decode string")
		panic(err)
	}

	encoded, err := x509.MarshalECPrivateKey(privKey)
	if err != nil {
		logger.Infof("unable to convert to bytes")
		panic(err)
	}

	httputil.NewResponse[any](c, http.MethodPost, gin.H{"str": encoded})
}
