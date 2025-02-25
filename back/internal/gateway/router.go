package gateway

import (
	"backend/internal/service"
	"backend/utils"
	"log"
	"os"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func InitRouter(services service.Service, basicMw *utils.BasicAuthMiddleware) {
	router := gin.Default()
	router.Use(cors.Default())

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
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "8081"
	}

	log.Printf("Run server on port = %s", port)
	router.Run(":" + port)
}
