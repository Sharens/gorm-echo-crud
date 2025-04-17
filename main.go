package main

import (
	"embed"
	"gorm-echo-crud/handler"
	"gorm-echo-crud/model"
	"html/template"
	"io"
	"io/fs"
	"net/http"
	"strings"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var templateFiles embed.FS

type TemplateRegistry struct {
	templates *template.Template
}

func (t *TemplateRegistry) Render(w io.Writer, name string, data interface{}, c echo.Context) error {
	return t.templates.ExecuteTemplate(w, name, data)
}

func main() {
	e := echo.New()

	db, err := gorm.Open(sqlite.Open("store.db"), &gorm.Config{})
	if err != nil {
		e.Logger.Fatal("Failed to connect database: ", err)
	}

	err = db.AutoMigrate(&model.Product{}, &model.Cart{}, &model.CartItem{}, &model.Category{})
	if err != nil {
		e.Logger.Fatal("Failed to migrate database: ", err)
	}
	e.Logger.Info("Database migration completed successfully.")

	templates := template.New("")

	err = fs.WalkDir(templateFiles, ".", func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if !d.IsDir() && strings.HasSuffix(path, ".html") {
			content, readErr := templateFiles.ReadFile(path)
			if readErr != nil {
				return readErr
			}
			_, parseErr := templates.Parse(string(content))
			if parseErr != nil {
				e.Logger.Errorf("Error parsing embedded template %s: %v", path, parseErr)
				return parseErr
			}
		}
		return nil
	})
	if err != nil {
		e.Logger.Fatalf("Error walking/parsing embedded templates: %v", err)
	}
	e.Renderer = &TemplateRegistry{templates: templates}

	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"*"},
		AllowMethods: []string{http.MethodGet, http.MethodPut, http.MethodPost, http.MethodDelete},
	}))

	productHandler := &handler.ProductHandler{DB: db}
	cartHandler := &handler.CartHandler{DB: db}
	categoryHandler := &handler.CategoryHandler{DB: db}

	e.GET("/", func(c echo.Context) error { return c.Render(http.StatusOK, "index.html", nil) })
	e.GET("/test-ui", func(c echo.Context) error { return c.Render(http.StatusOK, "test-ui.html", nil) })

	productGroup := e.Group("/products")
	productGroup.POST("", productHandler.CreateProduct)
	productGroup.GET("", productHandler.GetProducts)
	productGroup.GET("/:id", productHandler.GetProduct)
	productGroup.PUT("/:id", productHandler.UpdateProduct)
	productGroup.DELETE("/:id", productHandler.DeleteProduct)

	cartGroup := e.Group("/cart")
	cartGroup.POST("", cartHandler.CreateCart)
	cartGroup.GET("/:cart_id", cartHandler.GetCart)
	cartGroup.DELETE("/:cart_id", cartHandler.DeleteCart)
	cartGroup.POST("/:cart_id/items", cartHandler.AddItemToCart)
	cartGroup.DELETE("/:cart_id/items/:item_id", cartHandler.RemoveItemFromCart)

	categoryGroup := e.Group("/categories")
	categoryGroup.POST("", categoryHandler.CreateCategory)
	categoryGroup.GET("", categoryHandler.GetCategories)
	categoryGroup.GET("/:id", categoryHandler.GetCategory)
	categoryGroup.PUT("/:id", categoryHandler.UpdateCategory)
	categoryGroup.DELETE("/:id", categoryHandler.DeleteCategory)

	e.Logger.Info("Starting server on http://localhost:1323")
	e.Logger.Fatal(e.Start(":1323"))
}
