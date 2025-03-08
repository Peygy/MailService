package service

import (
	"backend/internal/model"
	"backend/utils"
	"encoding/base64"
	"encoding/json"
	"log"
	"net/http"
	"slices"
	"strings"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
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

func (ms *mailService) GetInboxMails(c *gin.Context) {
	userID := c.MustGet("userID").(uint)

	var user model.User
	if err := ms.db.Where("id = ?", userID).First(&user).Error; err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "Invalid email"})
		return
	}

	var mails []model.Mail
	err := ms.db.Find(&mails).Error
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Error fetching mails"})
		return
	}

	newMails := make([]model.Mail, 0, len(mails))
	for _, mail := range mails {
		var receivers map[string]interface{}
		if err := json.Unmarshal(mail.Receivers.Bytes, &receivers); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"message": "Error decoding receivers 1"})
			return
		}

		decodedBytes, err := base64.StdEncoding.DecodeString(receivers["Bytes"].(string))
		if err != nil {
			log.Println("Error decoding Base64:", err)
			c.JSON(http.StatusInternalServerError, gin.H{"message": "Error decoding Base64"})
			return
		}

		var recs1 []string
		var recs2 string
		if err := json.Unmarshal(decodedBytes, &recs1); err != nil {
			if err := json.Unmarshal(decodedBytes, &recs2); err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"message": "Error decoding receivers 2"})
				return
			}
		}

		if slices.Contains(recs1, user.Email) || strings.Contains(recs2, user.Email) {
			newMails = append(newMails, mail)
		}
	}

	c.JSON(http.StatusOK, gin.H{"mails": newMails})
}

func (ms *mailService) GetSentMails(c *gin.Context) {
	userID := c.MustGet("userID").(uint)

	var user model.User
	if err := ms.db.Where("id = ?", userID).First(&user).Error; err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "Invalid credentials"})
		return
	}

	var mails []model.Mail
	err := ms.db.Where("sender = ? AND is_deleted = false", user.Email).Find(&mails).Error
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Error fetching sent mails"})
		return
	}

	responseMails := make([]map[string]interface{}, 0, len(mails))
	for _, mail := range mails {
		var receivers map[string]interface{}
		if err := json.Unmarshal(mail.Receivers.Bytes, &receivers); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"message": "Error decoding receivers"})
			return
		}

		decodedBytes, err := base64.StdEncoding.DecodeString(receivers["Bytes"].(string))
		if err != nil {
			log.Println("Error decoding Base64:", err)
			c.JSON(http.StatusInternalServerError, gin.H{"message": "Error decoding Base64"})
			return
		}

		var recs []string
		if err := json.Unmarshal(decodedBytes, &recs); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"message": "Error decoding receivers"})
			return
		}

		responseMails = append(responseMails, map[string]interface{}{
			"ID":        mail.ID,
			"Sender":    mail.Sender,
			"Receivers": string(decodedBytes),
			"Subject":   mail.Subject,
			"Body":      mail.Body,
			"CreatedAt": mail.CreatedAt,
		})
	}

	c.JSON(http.StatusOK, gin.H{"mails": responseMails})
}

func (ms *mailService) SendMail(c *gin.Context) {
	userID := c.MustGet("userID").(uint)

	var mailData struct {
		Receivers []string `json:"receivers"`
		Subject   string   `json:"subject"`
		Body      string   `json:"body"`
	}

	if err := c.ShouldBindJSON(&mailData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid input"})
		return
	}

	var user model.User
	if err := ms.db.Where("id = ?", userID).First(&user).Error; err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "Invalid credentials"})
		return
	}

	mail := model.Mail{
		Sender:  user.Email,
		Subject: mailData.Subject,
		Body:    mailData.Body,
	}
	mail.Receivers.Set(mailData.Receivers)

	if err := utils.SendMailSMTP(mail, mailData); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Error sending email through SMTP"})
		return
	}

	if err := ms.db.Create(&mail).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Error sending mail"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{})
}

func (ms *mailService) ClearTrash(c *gin.Context) {
	userID := c.MustGet("userID").(uint)

	var user model.User
	if err := ms.db.Where("id = ?", userID).First(&user).Error; err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "Invalid credentials"})
		return
	}

	err := ms.db.Where("sender = ? AND is_deleted = true", user.Email).Delete(&model.Mail{}).Error
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Error clearing trash"})
		return
	}

	c.JSON(http.StatusOK, gin.H{})
}
