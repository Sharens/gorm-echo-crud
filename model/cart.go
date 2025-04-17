package model

import "gorm.io/gorm"

type Cart struct {
	gorm.Model
	CartItems []CartItem `json:"items" gorm:"foreignKey:CartID"`
}
