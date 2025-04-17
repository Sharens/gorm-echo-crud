package model

import "gorm.io/gorm"

type Category struct {
	gorm.Model
	Name     string    `json:"name" gorm:"unique;not null"`
	Products []Product `json:"products,omitempty" gorm:"foreignKey:CategoryID"`
}
