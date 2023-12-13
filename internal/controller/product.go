package controller

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/julianstephens/license-server/internal/controller/httputil"
	"github.com/julianstephens/license-server/internal/licensemanager"
	"github.com/julianstephens/license-server/internal/model"
	"github.com/julianstephens/license-server/service"
	"github.com/mitchellh/mapstructure"
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
	res, err := service.GetAll[model.Product](base.DB)
	if err != nil {
		httputil.NewError(c, http.StatusInternalServerError, err)
		return
	}

	httputil.NewResponse(c, http.MethodGet, res)
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

	res, err := service.FindById[model.Product](base.DB, productId)
	if err != nil {
		httputil.NewError(c, http.StatusInternalServerError, err)
		return
	}

	httputil.NewResponse(c, http.MethodGet, res)
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

	httputil.NewResponse(c, http.MethodPost, res)
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

	httputil.NewResponse(c, http.MethodPut, res)
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

	res, err := service.Delete(base.DB, productId, &model.Product{})
	if err != nil {
		httputil.NewError(c, http.StatusInternalServerError, err)
		return
	}

	httputil.NewResponse(c, http.MethodDelete, res)
}

func (base *Controller) CreateProductKeyPair(c *gin.Context) {
	lm := licensemanager.LicenseManager{Config: base.Config}

	productId, err := service.GetId(c)
	if err != nil {
		httputil.NewError(c, http.StatusBadRequest, err)
		return
	}

	product, err := service.FindById[model.Product](base.DB, productId)
	if err != nil {
		httputil.NewError(c, http.StatusInternalServerError, err)
		return
	}

	kp, err := lm.CreateProductKeyPair(product.Name, productId)
	if err != nil {
		httputil.NewError(c, http.StatusInternalServerError, err)
		return
	}

	httputil.NewResponse(c, http.MethodDelete, kp)
}
