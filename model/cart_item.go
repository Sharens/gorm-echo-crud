package model

import "gorm.io/gorm"

type CartItem struct {
	gorm.Model
	CartID    uint    `json:"-" gorm:"not null;index"`
	ProductID uint    `json:"product_id" gorm:"not null;index"`
	Quantity  int     `json:"quantity" gorm:"not null;check:quantity > 0"`
	Product   Product `json:"product" gorm:"foreignKey:ProductID;references:ID"`
}
