package main

import (
	"backend/cmd"
	"backend/internal/model"
	"backend/utils"
	"flag"
	"log"

	"github.com/joho/godotenv"
	"golang.org/x/crypto/bcrypt"
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

	devFlag := flag.Bool("dev", false, "Run in development mode")
	flag.Parse()

	if *devFlag {
		devRun()
	} else {
		db.AutoMigrate(&model.User{}, &model.Mail{})
		log.Println("Database migration completed!")
	}
}

func main() {
	app := cmd.NewApp(db)
	app.Run()
}

func devRun() {
	if err := db.Migrator().DropTable(&model.User{}, &model.Mail{}); err != nil {
		log.Fatal("Failed to drop tables:", err)
	}
	if err := db.AutoMigrate(&model.User{}, &model.Mail{}); err != nil {
		log.Fatal("Failed to migrate tables:", err)
	}

	users := []model.User{
		{Email: "test1@gomail.kurs", Password: "12344", Role: model.RoleUser},
		{Email: "test2@gomail.kurs", Password: "12344", Role: model.RoleUser},
		{Email: "test3@gomail.kurs", Password: "12344", Role: model.RoleUser},
		{Email: "admin@gomail.kurs", Password: "12344adm", Role: model.RoleAdmin},
	}

	for _, user := range users {
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
		if err != nil {
			log.Printf("Failed to hash password for user %s: %v", user.Email, err)
			continue
		}
		user.Password = string(hashedPassword)

		if err := db.Create(&user).Error; err != nil {
			log.Printf("Failed to create user %s: %v", user.Email, err)
		} else {
			log.Printf("User %s created successfully", user.Email)
		}
	}
}
