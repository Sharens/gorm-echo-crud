package main

import (
	// "gorm.io/gorm"
	// "gorm.io/driver/sqlite"
	"net/http"
	// "log"
	"github.com/labstack/echo/v4"
)

func main() {
	e := echo.New()
	e.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "Hello World")
	})

	e.Logger.Fatal(e.Start(":1323"))
}
