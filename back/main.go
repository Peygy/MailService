package main

import (
	"backend/cmd"
	"backend/internal/model"
	"backend/utils"
	"log"

	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var (
	db *gorm.DB
)

func init() {
	if err := godotenv.Load(); err != nil {
		log.Print("No .env file found")
	}

	var err error
	db, err = gorm.Open(postgres.Open(utils.GetEnv("DB_CONF", "")), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect to the database:", err)
	}

	db.AutoMigrate(&model.User{}, &model.Mail{})
	log.Println("Database migration completed!")
}

func main() {
	app := cmd.NewApp(db)
	app.Run()
}
