package models

import "gorm.io/gorm"

type User struct {
	gorm.Model
	Email     string `json:"email" binding:"required" gorm:"unique"`
	FirstName string `json:"firstName" binding:"required"`
	LastName  string `json:"lastName" binding:"required"`
	Password  string `json:"password" binding:"required"`
}
