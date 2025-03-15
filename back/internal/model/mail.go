package model

import (
	"github.com/jackc/pgx/pgtype"
	"github.com/jinzhu/gorm"
	"github.com/lib/pq"
)

type (
	Mail struct {
		gorm.Model
		Sender    string       `gorm:"not null"`
		Receivers pgtype.JSONB `gorm:"type:jsonb;default:'[]';not null"`

		Subject string
		Body    string
	}

	Trash struct {
		gorm.Model
		UserId   uint          `gorm:"uniqueIndex;not null"`
		Archived pq.Int64Array `gorm:"type:integer[]"`
		Deleted  pq.Int64Array `gorm:"type:integer[]"`
	}
)
