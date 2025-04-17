package handler

import (
	"errors"
	"gorm-echo-crud/model"
	"net/http"
	"strconv"
	"strings"

	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

type CategoryHandler struct {
	DB *gorm.DB
}

func (h *CategoryHandler) CreateCategory(c echo.Context) error {
	cat := new(model.Category)
	if err := c.Bind(cat); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid input: "+err.Error())
	}

	if cat.Name == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid input: Name is required")
	}

	result := h.DB.Create(&cat)
	if result.Error != nil {

		if strings.Contains(result.Error.Error(), "UNIQUE constraint failed") {
			return echo.NewHTTPError(http.StatusConflict, "Category name already exists")
		}
		return echo.NewHTTPError(http.StatusInternalServerError, "Could not create category: "+result.Error.Error())
	}

	return c.JSON(http.StatusCreated, cat)
}

func (h *CategoryHandler) GetCategories(c echo.Context) error {
	var categories []model.Category

	result := h.DB.Find(&categories)

	if result.Error != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Could not retrieve categories: "+result.Error.Error())
	}

	return c.JSON(http.StatusOK, categories)
}

func (h *CategoryHandler) GetCategory(c echo.Context) error {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil || id <= 0 {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid ID format")
	}

	var category model.Category

	result := h.DB.First(&category, id)

	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return echo.NewHTTPError(http.StatusNotFound, "Category not found")
		}
		return echo.NewHTTPError(http.StatusInternalServerError, "Database error: "+result.Error.Error())
	}

	return c.JSON(http.StatusOK, category)
}

func (h *CategoryHandler) UpdateCategory(c echo.Context) error {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil || id <= 0 {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid ID format")
	}

	var existingCategory model.Category
	findResult := h.DB.First(&existingCategory, id)
	if findResult.Error != nil {
		if errors.Is(findResult.Error, gorm.ErrRecordNotFound) {
			return echo.NewHTTPError(http.StatusNotFound, "Category not found")
		}
		return echo.NewHTTPError(http.StatusInternalServerError, "Database error finding category: "+findResult.Error.Error())
	}

	input := make(map[string]interface{})
	if err := c.Bind(&input); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid input: "+err.Error())
	}

	if newName, ok := input["name"].(string); ok {
		if newName == "" {
			return echo.NewHTTPError(http.StatusBadRequest, "Invalid input: Name cannot be empty")
		}

		if newName != existingCategory.Name {
			existingCategory.Name = newName
		} else {

			return c.JSON(http.StatusOK, existingCategory)
		}
	} else {

		return echo.NewHTTPError(http.StatusBadRequest, "Invalid input: 'name' field is required for update")
	}

	saveResult := h.DB.Save(&existingCategory)
	if saveResult.Error != nil {
		if strings.Contains(saveResult.Error.Error(), "UNIQUE constraint failed") {
			return echo.NewHTTPError(http.StatusConflict, "Category name already exists")
		}
		return echo.NewHTTPError(http.StatusInternalServerError, "Could not update category: "+saveResult.Error.Error())
	}

	return c.JSON(http.StatusOK, existingCategory)
}

func (h *CategoryHandler) DeleteCategory(c echo.Context) error {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil || id <= 0 {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid ID format")
	}

	var productCount int64
	h.DB.Model(&model.Product{}).Where("category_id = ?", id).Count(&productCount)
	if productCount > 0 {
		return echo.NewHTTPError(http.StatusConflict, "Cannot delete category: contains existing products")
	}

	result := h.DB.Delete(&model.Category{}, id)

	if result.Error != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Could not delete category: "+result.Error.Error())
	}

	if result.RowsAffected == 0 {
		return echo.NewHTTPError(http.StatusNotFound, "Category not found or already deleted")
	}

	return c.NoContent(http.StatusNoContent)
}
