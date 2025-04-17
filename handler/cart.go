package handler

import (
	"errors"
	"gorm-echo-crud/model"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

type CartHandler struct {
	DB *gorm.DB
}

func (h *CartHandler) CreateCart(c echo.Context) error {
	cart := model.Cart{}
	result := h.DB.Create(&cart)
	if result.Error != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Could not create cart: "+result.Error.Error())
	}
	return c.JSON(http.StatusCreated, cart)
}

func (h *CartHandler) GetCart(c echo.Context) error {
	cartIDStr := c.Param("cart_id")
	cartID, err := strconv.ParseUint(cartIDStr, 10, 32)
	if err != nil || cartID == 0 {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid Cart ID format")
	}

	var cart model.Cart

	result := h.DB.Scopes(model.CartWithDetails).First(&cart, uint(cartID))

	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return echo.NewHTTPError(http.StatusNotFound, "Cart not found")
		}
		return echo.NewHTTPError(http.StatusInternalServerError, "Database error: "+result.Error.Error())
	}

	return c.JSON(http.StatusOK, cart)
}

type AddItemRequest struct {
	ProductID uint `json:"product_id"`
	Quantity  int  `json:"quantity"`
}

func (h *CartHandler) AddItemToCart(c echo.Context) error {
	cartIDStr := c.Param("cart_id")
	cartID, err := strconv.ParseUint(cartIDStr, 10, 32)
	if err != nil || cartID == 0 {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid Cart ID format")
	}

	req := new(AddItemRequest)
	if err := c.Bind(req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid request body: "+err.Error())
	}

	if req.ProductID == 0 || req.Quantity <= 0 {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid input: ProductID and positive Quantity are required")
	}

	var cart model.Cart
	if err := h.DB.First(&cart, uint(cartID)).Error; err != nil {
	}
	var product model.Product
	if err := h.DB.First(&product, req.ProductID).Error; err != nil {
	}

	var existingItem model.CartItem
	err = h.DB.Scopes(model.CartItemForCart(uint(cartID))).Where("product_id = ?", req.ProductID).First(&existingItem).Error

	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return echo.NewHTTPError(http.StatusInternalServerError, "Database error checking existing item: "+err.Error())
	}

	var itemToReturn model.CartItem
	var httpStatus = http.StatusOK

	if errors.Is(err, gorm.ErrRecordNotFound) {

		newItem := model.CartItem{CartID: uint(cartID), ProductID: req.ProductID, Quantity: req.Quantity}
		result := h.DB.Create(&newItem)
		if result.Error != nil {
		}
		itemToReturn = newItem
		httpStatus = http.StatusCreated
	} else {

		existingItem.Quantity += req.Quantity
		if existingItem.Quantity <= 0 {
			delResult := h.DB.Delete(&existingItem)
			if delResult.Error != nil {
			}
			return c.NoContent(http.StatusNoContent)
		}
		result := h.DB.Save(&existingItem)
		if result.Error != nil {
		}
		itemToReturn = existingItem
	}

	h.DB.Scopes(model.CartItemWithProduct).First(&itemToReturn, itemToReturn.ID)

	return c.JSON(httpStatus, itemToReturn)
}

func (h *CartHandler) RemoveItemFromCart(c echo.Context) error {
	cartIDStr := c.Param("cart_id")
	cartID, err := strconv.ParseUint(cartIDStr, 10, 32)
	if err != nil || cartID == 0 {
	}

	itemIDStr := c.Param("item_id")
	itemID, err := strconv.ParseUint(itemIDStr, 10, 32)
	if err != nil || itemID == 0 {
	}

	result := h.DB.Scopes(model.CartItemForCart(uint(cartID))).Delete(&model.CartItem{}, uint(itemID))

	if result.Error != nil {
	}
	if result.RowsAffected == 0 {
		return echo.NewHTTPError(http.StatusNotFound, "Cart item not found in the specified cart")
	}

	return c.NoContent(http.StatusNoContent)
}

func (h *CartHandler) DeleteCart(c echo.Context) error {
	cartIDStr := c.Param("cart_id")
	cartID, err := strconv.ParseUint(cartIDStr, 10, 32)
	if err != nil || cartID == 0 {
	}

	err = h.DB.Transaction(func(tx *gorm.DB) error {

		if err := tx.Scopes(model.CartItemForCart(uint(cartID))).Delete(&model.CartItem{}).Error; err != nil {
			return err
		}

		result := tx.Delete(&model.Cart{}, uint(cartID))
		if result.Error != nil {
			return result.Error
		}
		if result.RowsAffected == 0 {
			return gorm.ErrRecordNotFound
		}

		return nil
	})

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return echo.NewHTTPError(http.StatusNotFound, "Cart not found")
		}
		return echo.NewHTTPError(http.StatusInternalServerError, "Could not delete cart: "+err.Error())
	}

	return c.NoContent(http.StatusNoContent)
}
