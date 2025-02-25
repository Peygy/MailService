package service

import (
	"backend/internal/model"
	"backend/utils"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
)

type (
	MailService interface {
		GetInboxMails(c *gin.Context)
		GetSentMails(c *gin.Context)
		SendMail(c *gin.Context)
		ClearTrash(c *gin.Context)
	}

	mailService struct {
		db *gorm.DB
	}
)

func NewMailService(db *gorm.DB) MailService {
	return &mailService{
		db: db,
	}
}

// Получение входящих писем
func (ms *mailService) GetInboxMails(c *gin.Context) {
	userID := c.MustGet("userID").(uint)

	var mails []model.Mail
	err := ms.db.Where("receiver = ? AND is_deleted = false", userID).Find(&mails).Error
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Error fetching inbox mails"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"mails": mails})
}

// Получение отправленных писем
func (ms *mailService) GetSentMails(c *gin.Context) {
	userID := c.MustGet("userID").(uint)

	var mails []model.Mail
	err := ms.db.Where("sender_id = ? AND is_deleted = false", userID).Find(&mails).Error
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Error fetching sent mails"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"mails": mails})
}

// Отправка письма
func (ms *mailService) SendMail(c *gin.Context) {
	userID := c.MustGet("userID").(uint)

	var mailData struct {
		Receiver string `json:"receiver"`
		Subject  string `json:"subject"`
		Body     string `json:"body"`
	}

	if err := c.ShouldBindJSON(&mailData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid input"})
		return
	}

	// Создаем новое письмо
	mail := model.Mail{
		SenderID: userID,
		Receiver: mailData.Receiver,
		Subject:  mailData.Subject,
		Body:     mailData.Body,
	}

	if err := ms.db.Create(&mail).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Error sending mail"})
		return
	}

	// Отправляем письмо через SMTP
	err := utils.SendMailSMTP(mailData.Receiver, mailData.Subject, mailData.Body)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Error sending email through SMTP"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "Mail sent successfully"})
}

// Очистка корзины
func (ms *mailService) ClearTrash(c *gin.Context) {
	userID := c.MustGet("userID").(uint)

	err := ms.db.Where("sender_id = ? AND is_deleted = true", userID).Delete(&model.Mail{}).Error
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Error clearing trash"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Trash cleared successfully"})
}
