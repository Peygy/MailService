package cmd

import (
	"backend/internal/gateway"
	"backend/internal/service"
	"backend/utils"
	"log"

	"github.com/jinzhu/gorm"
	_ "github.com/lib/pq"
)

type App struct {
	db *gorm.DB
}

func NewApp(db *gorm.DB) *App {
	return &App{
		db: db,
	}
}

func (a *App) Run() {
	mailServ := service.NewMailService(a.db)
	authServ := service.NewAuthService(a.db)
	services := service.Service{
		MailService: mailServ,
		AuthService: authServ,
	}

	basicAuthMw := utils.Init(a.db)

	log.Println("Initialize router")
	gateway.InitRouter(services, basicAuthMw)
}
