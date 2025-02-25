package main

import (
	"backend/cmd"
	"backend/internal/model"
	"log"

	"github.com/jinzhu/gorm"
	_ "github.com/lib/pq"
)

var (
	db *gorm.DB
)

func init() {
	dsn := "user=postgres password=postgres dbname=mails port=5432 sslmode=disable"

	db, err := gorm.Open("postgres", dsn)
	if err != nil {
		log.Fatal("Failed to connect to the database:", err)
	}

	db.AutoMigrate(&model.User{}, &model.Mail{})
}

func main() {
	app := cmd.NewApp(db)
	app.Run()
}
