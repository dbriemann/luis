package models

import "gorm.io/gorm"

type User struct {
	gorm.Model

	Email   string `gorm:"uniqueIndex"`
	Secret  string
	Name    string
	IsAdmin bool
}
