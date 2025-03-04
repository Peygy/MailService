package model

const (
	RoleUser  = "user"
	RoleAdmin = "admin"
)

type User struct {
	Id       uint   `gorm:"primaryKey"`
	Email    string `gorm:"uniqueIndex;not null"`
	Password string `gorm:"not null"`
	Role     string `gorm:"type:varchar(10);not null;default:'user'"`
}
