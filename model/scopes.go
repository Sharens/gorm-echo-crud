package model

import "gorm.io/gorm"

func ProductWithCategory(db *gorm.DB) *gorm.DB {
	return db.Preload("Category")
}

func CartWithDetails(db *gorm.DB) *gorm.DB {

	return db.Preload("CartItems.Product")
}

func CartItemForCart(cartID uint) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		return db.Where("cart_id = ?", cartID)
	}
}

func CartItemWithProduct(db *gorm.DB) *gorm.DB {
	return db.Preload("Product")
}
