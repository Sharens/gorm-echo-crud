package main

import (
	// "gorm.io/gorm"
	// "gorm.io/driver/sqlite"

	"fmt"
	"html/template"
	"io"
	"net/http"
	"strings"

	// "log"
	"gorm-echo-crud/handler"
	product "gorm-echo-crud/handler"

	"github.com/labstack/echo/v4"
)

type TemplateRenderer struct {
	templates *template.Template
}

func (t *TemplateRenderer) Render(w io.Writer, name string, data interface{}, c echo.Context) error {
	templ := t.templates.Lookup(name)
	if templ == nil {
		return fmt.Errorf("template %s not found", name)
	}

	return templ.ExecuteTemplate(w, name, data)
}

func main() {
	e := echo.New()
	result := strings.Join(product.GetProductList(), ", ")
	templates, err := template.ParseGlob("view/*.html")
	if err != nil {
		e.Logger.Fatal("Error loading templates: ", err)
	}

	renderer := &TemplateRenderer{
		templates: templates,
	}

	e.Renderer = renderer

	e.GET("/", func(c echo.Context) error {
		return c.Render(http.StatusOK, "index.html", nil)
	})

	e.GET("/products", func(c echo.Context) error {
		return c.String(http.StatusOK, result)
	})

	handler.RegisterRoutes(e)

	e.Logger.Fatal(e.Start(":1323"))
}
