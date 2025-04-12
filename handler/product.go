package handler

import (
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
)

var productList = []string{"laptop", "komputer", "mikrofon"}

func GetProductList() []string {
	return productList
}

func GetProduct(c echo.Context) error {
	idStr := c.Param("id")
	idInt, err := strconv.Atoi(idStr)

	if err != nil {
		return c.String(http.StatusBadRequest, "Nieprawidłowy format ID: "+idStr)
	}

	return c.String(http.StatusOK, productList[idInt-1])
}

func AddProduct(c echo.Context) error {
	product_name := c.FormValue("product_name")

	return c.String(http.StatusOK, "product_name:"+product_name)
}

// func (h *Handler) GetArticle(c echo.Context) error {
// 	slug := c.Param("slug")
// 	a, err := h.articleStore.GetBySlug(slug)

// 	if err != nil {
// 		return c.JSON(http.StatusInternalServerError, utils.NewError(err))
// 	}

// 	if a == nil {
// 		return c.JSON(http.StatusNotFound, utils.NotFound())
// 	}

// 	return c.JSON(http.StatusOK, newArticleResponse(c, a))
// }
