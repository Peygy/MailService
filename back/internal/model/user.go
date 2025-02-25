package model

import "github.com/jinzhu/gorm"

type User struct {
	gorm.Model
	Email string `gorm:"unique"`
	Password string
	RoleID   uint
	Role     Role
}

type Role struct {
	gorm.Model
	Name string `gorm:"unique"`
}
