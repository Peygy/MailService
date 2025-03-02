package model

import (
	"github.com/jackc/pgx/pgtype"
	"github.com/jinzhu/gorm"
)

type Mail struct {
	gorm.Model
	Sender    string       `gorm:"not null"`
	Receivers pgtype.JSONB `gorm:"type:jsonb;default:'[]';not null"`

	Subject string
	Body    string

	IsRead    bool `gorm:"default:false"`
	IsDeleted bool `gorm:"default:false"`
}
