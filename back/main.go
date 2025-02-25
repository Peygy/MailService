package main

import (
	"backend/cmd"
	"backend/internal/model"
	"log"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var (
	db *gorm.DB
)

func init() {
	dsn := "user=postgres password=postgres dbname=mails port=5432 sslmode=disable"

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect to the database:", err)
	}

	db.AutoMigrate(&model.User{}, &model.Mail{})
}

func main() {
	app := cmd.NewApp(db)
	app.Run()
}
