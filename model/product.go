package model

import "gorm.io/gorm"

type Product struct {
	gorm.Model
	Name       string   `json:"name" gorm:"not null"`
	Price      float64  `json:"price" gorm:"not null;check:price > 0"`
	CategoryID uint     `json:"category_id" gorm:"index;not null"`
	Category   Category `json:"category,omitempty" gorm:"foreignKey:CategoryID"`
}
