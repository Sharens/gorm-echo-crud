package handler

import (
	"errors"
	"gorm-echo-crud/model"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

type ProductHandler struct {
	DB *gorm.DB
}

func (h *ProductHandler) CreateProduct(c echo.Context) error {
	prod := new(model.Product)
	if err := c.Bind(prod); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid input: "+err.Error())
	}

	if prod.Name == "" || prod.Price <= 0 {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid input: Name cannot be empty and Price must be positive")
	}

	result := h.DB.Create(&prod)
	if result.Error != nil {

		return echo.NewHTTPError(http.StatusInternalServerError, "Could not create product: "+result.Error.Error())
	}

	return c.JSON(http.StatusCreated, prod)
}

func (h *ProductHandler) GetProducts(c echo.Context) error {
	var products []model.Product

	result := h.DB.Find(&products)
	if result.Error != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Could not retrieve products: "+result.Error.Error())
	}

	return c.JSON(http.StatusOK, products)
}

func (h *ProductHandler) GetProduct(c echo.Context) error {
	idStr := c.Param("id")

	id, err := strconv.Atoi(idStr)
	if err != nil || id <= 0 {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid ID format")
	}

	var product model.Product

	result := h.DB.First(&product, id)

	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return echo.NewHTTPError(http.StatusNotFound, "Product not found")
		}
		return echo.NewHTTPError(http.StatusInternalServerError, "Database error: "+result.Error.Error())
	}

	return c.JSON(http.StatusOK, product)
}

func (h *ProductHandler) UpdateProduct(c echo.Context) error {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil || id <= 0 {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid ID format")
	}

	var existingProduct model.Product
	findResult := h.DB.First(&existingProduct, id)
	if findResult.Error != nil {
		if errors.Is(findResult.Error, gorm.ErrRecordNotFound) {
			return echo.NewHTTPError(http.StatusNotFound, "Product not found")
		}
		return echo.NewHTTPError(http.StatusInternalServerError, "Database error: "+findResult.Error.Error())
	}

	updatedData := new(model.Product)
	if err := c.Bind(updatedData); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid input: "+err.Error())
	}

	if updatedData.Name == "" || updatedData.Price <= 0 {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid input: Name cannot be empty and Price must be positive")
	}

	existingProduct.Name = updatedData.Name
	existingProduct.Price = updatedData.Price

	saveResult := h.DB.Save(&existingProduct)
	if saveResult.Error != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Could not update product: "+saveResult.Error.Error())
	}

	return c.JSON(http.StatusOK, existingProduct)
}

func (h *ProductHandler) DeleteProduct(c echo.Context) error {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil || id <= 0 {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid ID format")
	}

	result := h.DB.Delete(&model.Product{}, id)

	if result.Error != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Could not delete product: "+result.Error.Error())
	}

	if result.RowsAffected == 0 {

		return echo.NewHTTPError(http.StatusNotFound, "Product not found or already deleted")
	}

	return c.NoContent(http.StatusNoContent)
}
