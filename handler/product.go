package handler

import (
	"gorm-echo-crud/model"
	"net/http"
	"strconv"
	"sync"

	"github.com/labstack/echo/v4"
)

var (
	products   = make(map[int]*model.Product)
	nextProdID = 1
	mu         sync.Mutex
)

func CreateProduct(c echo.Context) error {
	prod := new(model.Product)
	if err := c.Bind(prod); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid input: "+err.Error())
	}

	if prod.Name == "" || prod.Price <= 0 {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid input: Name cannot be empty and Price must be positive")
	}

	mu.Lock()
	prod.ID = nextProdID
	products[prod.ID] = prod
	nextProdID++
	mu.Unlock()

	return c.JSON(http.StatusCreated, prod)
}

func GetProducts(c echo.Context) error {
	mu.Lock()
	defer mu.Unlock()

	productList := make([]*model.Product, 0, len(products))
	for _, prod := range products {
		productList = append(productList, prod)
	}

	return c.JSON(http.StatusOK, productList)
}

func GetProduct(c echo.Context) error {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid ID format")
	}

	mu.Lock()
	prod, exists := products[id]
	mu.Unlock()

	if !exists {
		return echo.NewHTTPError(http.StatusNotFound, "Product not found")
	}

	return c.JSON(http.StatusOK, prod)
}

func UpdateProduct(c echo.Context) error {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid ID format")
	}

	mu.Lock()
	existingProd, exists := products[id]
	if !exists {
		mu.Unlock()
		return echo.NewHTTPError(http.StatusNotFound, "Product not found")
	}
	mu.Unlock()

	updatedProdData := new(model.Product)
	if err := c.Bind(updatedProdData); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid input: "+err.Error())
	}

	if updatedProdData.Name == "" || updatedProdData.Price <= 0 {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid input: Name cannot be empty and Price must be positive")
	}

	mu.Lock()
	existingProd.Name = updatedProdData.Name
	existingProd.Price = updatedProdData.Price

	mu.Unlock()

	return c.JSON(http.StatusOK, existingProd)
}

func DeleteProduct(c echo.Context) error {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid ID format")
	}

	mu.Lock()
	_, exists := products[id]
	if !exists {
		mu.Unlock()
		return echo.NewHTTPError(http.StatusNotFound, "Product not found")
	}
	delete(products, id)
	mu.Unlock()

	return c.NoContent(http.StatusNoContent)
}
