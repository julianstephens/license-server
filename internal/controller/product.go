package controller

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/mitchellh/mapstructure"

	"github.com/julianstephens/license-server/pkg/httputil"
	"github.com/julianstephens/license-server/pkg/licensemanager"
	"github.com/julianstephens/license-server/pkg/logger"
	"github.com/julianstephens/license-server/pkg/model"
	"github.com/julianstephens/license-server/pkg/service"
)

// GetProducts godoc
// @Summary Get all products
// @Description retrieves all products
// @Tags products
// @Security ApiKey
// @Success 200 {object} httputil.HTTPResponse[[]model.Product]
// @Failure 500 {object} httputil.HTTPError
// @Router /admin/products [get]
func (base *Controller) GetProducts(c *gin.Context) {
	logger.Infof("retrieving all products as user <%s>", c.GetString("user"))

	res, err := service.GetAll[model.Product](base.DB)
	if err != nil {
		httputil.NewError(c, http.StatusInternalServerError, err)
		return
	}

	httputil.NewResponse(c, res, httputil.Options{IsCrudHandler: true, HttpMsgMethod: httputil.Get})
}

// GetProduct godoc
// @Summary Get a product
// @Description retrieve a specific product
// @Tags products
// @Security ApiKey
// @Success 200 {object} httputil.HTTPResponse[model.Product]
// @Failure 400 {object} httputil.HTTPError
// @Failure 500 {object} httputil.HTTPError
// @Router /admin/products/:id [get]
func (base *Controller) GetProduct(c *gin.Context) {
	productId, err := service.GetId(c)
	if err != nil {
		httputil.NewError(c, http.StatusBadRequest, err)
		return
	}
	logger.Infof("retrieving product <%s> as user <%s>", productId, c.GetString("user"))

	res, err := service.FindById[model.Product](base.DB, productId)
	if err != nil {
		httputil.NewError(c, http.StatusInternalServerError, err)
		return
	}

	httputil.NewResponse(c, res, httputil.Options{IsCrudHandler: true, HttpMsgMethod: httputil.Get})
}

// AddProduct godoc
// @Summary Add a product
// @Description creates a new product
// @Tags products
// @Security ApiKey
// @Param data body model.Product true "new product info"
// @Success 200 {object} httputil.HTTPResponse[model.Product]
// @Failure 400 {object} httputil.HTTPError
// @Failure 500 {object} httputil.HTTPError
// @Router /admin/products [post]
func (base *Controller) AddProduct(c *gin.Context) {
	logger.Infof("creating new product as user <%s>", c.GetString("user"))
	var product *model.Product

	if err := c.ShouldBindJSON(&product); err != nil {
		httputil.HandleFieldError(c, err)
		return
	}

	res, err := service.Create(base.DB, product)
	if err != nil {
		httputil.NewError(c, http.StatusInternalServerError, err)
		return
	}

	httputil.NewResponse(c, res, httputil.Options{IsCrudHandler: true, HttpMsgMethod: httputil.Post})
}

// UpdateProduct godoc
// @Summary Update a product
// @Description updates a specific product
// @Tags products
// @Security ApiKey
// @Param data body model.Product true "updated product info"
// @Success 200 {object} httputil.HTTPResponse[model.Product]
// @Failure 400 {object} httputil.HTTPError
// @Failure 500 {object} httputil.HTTPError
// @Router /admin/products/:id [put]
func (base *Controller) UpdateProduct(c *gin.Context) {
	var productUpdates *model.Product

	productId, err := service.GetId(c)
	if err != nil {
		httputil.NewError(c, http.StatusBadRequest, err)
		return
	}
	logger.Infof("updating product <%s> as user <%s>", productId, c.GetString("user"))

	if err := c.ShouldBindJSON(&productUpdates); err != nil {
		httputil.HandleFieldError(c, err)
		return
	}

	var decoded map[string]any
	if err := mapstructure.Decode(productUpdates, decoded); err != nil {
		httputil.NewError(c, http.StatusInternalServerError, err)
		return
	}

	res, err := service.Update[model.Product](base.DB, productId, decoded)
	if err != nil {
		httputil.NewError(c, http.StatusInternalServerError, err)
		return
	}

	httputil.NewResponse(c, res, httputil.Options{IsCrudHandler: true, HttpMsgMethod: httputil.Put})
}

// DeleteProduct godoc
// @Summary Delete a product
// @Description deletes a specific product
// @Tags products
// @Security ApiKey
// @Success 200 {object} httputil.HTTPResponse[model.Product]
// @Failure 400 {object} httputil.HTTPError
// @Failure 500 {object} httputil.HTTPError
// @Router /admin/products/:id [delete]
func (base *Controller) DeleteProduct(c *gin.Context) {
	productId, err := service.GetId(c)
	if err != nil {
		httputil.NewError(c, http.StatusBadRequest, err)
		return
	}
	logger.Infof("deleting product <%s> as user <%s>", productId, c.GetString("user"))

	res, err := service.Delete(base.DB, productId, &model.Product{})
	if err != nil {
		httputil.NewError(c, http.StatusInternalServerError, err)
		return
	}

	httputil.NewResponse(c, res, httputil.Options{IsCrudHandler: true, HttpMsgMethod: httputil.Delete})
}

// CreateProductKeyPair godoc
// @Summary Add a product key pair
// @Description creates an ed25519 key pair for a specific product and version
// @Tags products
// @Security ApiKey
// @Success 200 {object} httputil.HTTPResponse[model.ProductKeyPair]
// @Failure 400 {object} httputil.HTTPError
// @Failure 500 {object} httputil.HTTPError
// @Router /admin/products/:id/key [get]
func (base *Controller) CreateProductKeyPair(c *gin.Context) {
	lm := licensemanager.LicenseManager{Config: base.Config, DB: base.DB}

	productId, err := service.GetId(c)
	if err != nil {
		httputil.NewError(c, http.StatusBadRequest, err)
		return
	}
	logger.Infof("generating new key pair for product <%s> as user <%s>", productId, c.GetString("user"))

	kp, err := lm.CreateProductKeyPair(productId)
	if err != nil {
		httputil.NewError(c, http.StatusInternalServerError, err)
		return
	}

	httputil.NewResponse(c, kp, httputil.Options{IsCrudHandler: true, HttpMsgMethod: httputil.Get})
}
