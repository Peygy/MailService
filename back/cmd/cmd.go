package cmd

import (
	"backend/internal/gateway"
	"backend/internal/model"
	"backend/internal/service"
	"backend/utils"
	"log"
	"time"

	_ "github.com/lib/pq"
	"gorm.io/gorm"
)

type App struct {
	db model.MailDB
}

func NewApp(db *gorm.DB) *App {
	return &App{
		db: model.NewMailDB(db),
	}
}

func (a *App) Run() {
	go func() {
		for {
			if err := utils.ReadMailIMAP(a.db); err != nil {
				log.Println("Error reading emails:", err)
			}
			time.Sleep(10 * time.Second)
		}
	}()

	mailServ := service.NewMailService(a.db)
	authServ := service.NewAuthService(a.db)
	adminServ := service.NewAdminService(a.db)

	services := service.Service{
		MailService:  mailServ,
		AuthService:  authServ,
		AdminService: adminServ,
	}

	basicAuthMw := utils.NewBasicAuthMiddleware(a.db)
	roleMw := utils.NewRoleMiddleware(a.db)

	log.Println("Initialize router")
	gateway.InitRouter(services, basicAuthMw, roleMw)
}
