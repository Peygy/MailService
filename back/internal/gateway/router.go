package gateway

import (
	"backend/internal/model"
	"backend/internal/service"
	"backend/utils"
	"log"
	"os"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func InitRouter(services service.Service, basicMw *utils.BasicAuthMiddleware,
	roleMw *utils.RoleMiddleware,
) {
	router := gin.Default()
	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Authorization", "Content-Type"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
	}))

	api := router.Group("/api/v1")
	{
		api.POST("/register", services.AuthService.RegisterUser)
		api.POST("/login", services.AuthService.Login)

		mail := api.Group("/mail", basicMw.Middleware())
		{
			mail.GET("/inbox", services.MailService.GetInboxMails)
			mail.GET("/sent", services.MailService.GetSentMails)
			mail.POST("/send", services.MailService.SendMail)
			mail.DELETE("/trash", services.MailService.ClearTrash)
		}

		admin := api.Group("/admin", basicMw.Middleware(), roleMw.Middleware(model.RoleAdmin))
		{
			admin.GET("/users", services.AdminService.GetAllUsers)
			admin.DELETE("/users/:id", services.AdminService.DeleteUser)
		}
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "8081"
	}

	log.Printf("Run server on port = %s", port)
	router.Run(":" + port)
}
