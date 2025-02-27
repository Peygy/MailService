package cmd

import (
	"backend/internal/gateway"
	"backend/internal/service"
	"backend/utils"
	"log"

	_ "github.com/lib/pq"
	"gorm.io/gorm"
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
	go func() {
		if err := utils.ReadMailIMAP(a.db); err != nil {
			log.Println("Error reading emails:", err)
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
