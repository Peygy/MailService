package model

import "github.com/jinzhu/gorm"

type Mail struct {
	gorm.Model
	SenderID  uint   `gorm:"not null"`
	Receiver  string `gorm:"not null"`
	Subject   string
	Body      string
	IsRead    bool `gorm:"default:false"`
	IsDeleted bool `gorm:"default:false"`
}
