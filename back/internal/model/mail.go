package model

import "github.com/jinzhu/gorm"

type Mail struct {
	gorm.Model
	SenderId   int    `gorm:"not null"`
	Sender     string `gorm:"not null"`
	ReceiverId int    `gorm:"not null"`
	Receiver   string `gorm:"not null"`

	Subject string
	Body    string

	IsRead    bool `gorm:"default:false"`
	IsDeleted bool `gorm:"default:false"`
}
