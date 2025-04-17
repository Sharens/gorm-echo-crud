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

	if prod.Name == "" || prod.Price <= 0 || prod.CategoryID == 0 {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid input: Name, positive Price, and CategoryID are required")
	}

	var category model.Category
	if err := h.DB.First(&category, prod.CategoryID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return echo.NewHTTPError(http.StatusBadRequest, "Invalid input: Category not found")
		}
		return echo.NewHTTPError(http.StatusInternalServerError, "Database error checking category: "+err.Error())
	}

	result := h.DB.Create(&prod)
	if result.Error != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Could not create product: "+result.Error.Error())
	}

	h.DB.Scopes(model.ProductWithCategory).First(&prod, prod.ID)

	return c.JSON(http.StatusCreated, prod)
}

func (h *ProductHandler) GetProducts(c echo.Context) error {
	var products []model.Product

	result := h.DB.Scopes(model.ProductWithCategory).Find(&products)
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

	result := h.DB.Scopes(model.ProductWithCategory).First(&product, id)

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
		return echo.NewHTTPError(http.StatusInternalServerError, "Database error finding product: "+findResult.Error.Error())
	}

	input := make(map[string]interface{})
	if err := c.Bind(&input); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid input: "+err.Error())
	}

	if newCatIDFloat, ok := input["category_id"].(float64); ok {
		newCatID := uint(newCatIDFloat)
		if newCatID != 0 && newCatID != existingProduct.CategoryID {
			var category model.Category
			if err := h.DB.First(&category, newCatID).Error; err != nil {
				if errors.Is(err, gorm.ErrRecordNotFound) {
					return echo.NewHTTPError(http.StatusBadRequest, "Invalid input: New category not found")
				}
				return echo.NewHTTPError(http.StatusInternalServerError, "Database error checking new category: "+err.Error())
			}
			input["category_id"] = newCatID
		} else if newCatID == 0 {
			delete(input, "category_id")
		} else {
			delete(input, "category_id")
		}
	}
	if name, ok := input["name"].(string); ok && name == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid input: Name cannot be empty")
	}
	if price, ok := input["price"].(float64); ok && price <= 0 {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid input: Price must be positive")
	}

	updateResult := h.DB.Model(&existingProduct).Updates(input)
	if updateResult.Error != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Could not update product: "+updateResult.Error.Error())
	}

	h.DB.Scopes(model.ProductWithCategory).First(&existingProduct, id)

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
