package main

import (
	// "gorm.io/gorm"
	// "gorm.io/driver/sqlite"
	"net/http"
	"strings"

	// "log"
	"gorm-echo-crud/handler"
	product "gorm-echo-crud/handler"

	"github.com/labstack/echo/v4"
)

func main() {
	e := echo.New()
	result := strings.Join(product.GetProductList(), ", ")

	e.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, result)
	})

	handler.RegisterRoutes(e)

	e.Logger.Fatal(e.Start(":1323"))
}
