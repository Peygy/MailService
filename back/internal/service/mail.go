package service

import (
	"backend/internal/model"
	"backend/utils"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"slices"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/kyroy/go-slices/int64s"
)

const (
	domain = "gomail.kurs"
)

type (
	MailService interface {
		GetInboxMails(c *gin.Context)
		GetSentMails(c *gin.Context)
		SendMail(c *gin.Context)
		GetTrash(c *gin.Context)
		UnArchiveMail(c *gin.Context)
		ArchiveMail(c *gin.Context)
		DeleteMail(c *gin.Context)
	}

	mailService struct {
		db model.MailDB
	}
)

func NewMailService(db model.MailDB) MailService {
	return &mailService{
		db: db,
	}
}

func (ms *mailService) GetInboxMails(c *gin.Context) {
	userID := c.MustGet("userID").(uint)

	var user model.User
	if err := ms.db.Where("id = ?", userID).First(&user).Error(); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "Invalid email"})
		return
	}

	var mails []model.Mail
	if err := ms.db.Find(&mails).Error(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Error fetching mails"})
		return
	}

	newMails := make([]model.Mail, 0, len(mails))
	for _, mail := range mails {
		if check, err := ms.checkEmailStat(userID, mail.ID); err != nil || !check {
			continue
		}

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
	if err := ms.db.Where("id = ?", userID).First(&user).Error(); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "Invalid credentials"})
		return
	}

	var mails []model.Mail
	if err := ms.db.Where("sender = ?", user.Email).Find(&mails).Error(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Error fetching sent mails"})
		return
	}

	responseMails := make([]map[string]interface{}, 0, len(mails))
	for _, mail := range mails {
		if check, err := ms.checkEmailStat(userID, mail.ID); err != nil || !check {
			continue
		}

		var receivers map[string]interface{}
		if err := json.Unmarshal(mail.Receivers.Bytes, &receivers); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"message": fmt.Sprintf("Error decoding receivers1: %v", err)})
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
			c.JSON(http.StatusInternalServerError, gin.H{"message": fmt.Sprintf("Error decoding receivers2: %v", err)})
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
	if err := ms.db.Where("id = ?", userID).First(&user).Error(); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "Invalid credentials"})
		return
	}

	mail := model.Mail{
		Sender:  user.Email,
		Subject: mailData.Subject,
		Body:    mailData.Body,
	}

	var filtered []string
	for _, rec := range mailData.Receivers {
		if !strings.Contains(rec, domain) {
			filtered = append(filtered, rec)
		}
	}

	if err := utils.SendMailSMTP(mail, filtered); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": fmt.Sprintf("Error sending email through SMTP: %v", err)})
		return
	}
	mail.Receivers.Set(mailData.Receivers)

	if err := ms.db.Create(&mail).Error(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Error sending mail"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{})
}

func (ms *mailService) GetTrash(c *gin.Context) {
	userID := c.MustGet("userID").(uint)

	type resp struct {
		ID        int
		Subject   string
		Body      string
		CreatedAt time.Time
	}
	var resps []resp

	var tr model.Trash
	if err := ms.db.Where("user_id = ?", userID).First(&tr).Error(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Archived mails not found"})
		return
	}

	for _, v := range tr.Archived {
		var mail model.Mail
		if err := ms.db.Where("id = ?", v).Find(&mail).Error(); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"message": "Error get mail for trash"})
			return
		}

		resps = append(resps, resp{ID: int(mail.ID), Subject: mail.Subject, Body: mail.Body, CreatedAt: mail.CreatedAt})
	}

	c.JSON(http.StatusOK, gin.H{"mails": resps})
}

func (ms *mailService) UnArchiveMail(c *gin.Context) {
	userID := c.MustGet("userID").(uint)
	mailID, _ := strconv.Atoi(c.Param("id"))

	var tr model.Trash
	if err := ms.db.Where("user_id = ?", userID).First(&tr).Error(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Archived mails not found"})
		return
	}

	mailIdx := int64s.IndexOf(tr.Archived, int64(mailID))
	if err := ms.db.Model(&model.Trash{}).
		Where("user_id = ?", userID).
		Update("archived", append(tr.Archived[:mailIdx], tr.Archived[mailIdx+1:]...)).Error(); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid mailID"})
		return
	}

	c.JSON(http.StatusOK, gin.H{})
}

func (ms *mailService) ArchiveMail(c *gin.Context) {
	userID := c.MustGet("userID").(uint)
	mailID, _ := strconv.Atoi(c.Param("id"))

	var tr model.Trash
	if err := ms.db.Where("user_id = ?", userID).First(&tr).Error(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Archived mails not found"})
		return
	}

	if err := ms.db.Model(&model.Trash{}).
		Where("user_id = ?", userID).
		Update("archived", append(tr.Archived, int64(mailID))).Error(); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid mailID"})
		return
	}

	c.JSON(http.StatusOK, gin.H{})
}

func (ms *mailService) DeleteMail(c *gin.Context) {
	userID := c.MustGet("userID").(uint)
	mailID, _ := strconv.Atoi(c.Param("id"))

	var tr model.Trash
	if err := ms.db.Where("user_id = ?", userID).First(&tr).Error(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Archived mails not found"})
		return
	}

	mailIdx := int64s.IndexOf(tr.Archived, int64(mailID))
	if err := ms.db.Model(&model.Trash{}).
		Where("user_id = ?", userID).
		Update("archived", append(tr.Archived[:mailIdx], tr.Archived[mailIdx+1:]...)).Error(); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid mailID"})
		return
	}

	if err := ms.db.Model(&model.Trash{}).
		Where("user_id = ?", userID).
		Update("deleted", append(tr.Deleted, int64(mailID))).Error(); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid mailID"})
		return
	}

	c.JSON(http.StatusOK, gin.H{})
}

func (ms *mailService) checkEmailStat(userID, mailID uint) (bool, error) {
	var tr model.Trash
	if err := ms.db.Model(&model.Trash{}).Where("user_id = ?", userID).First(&tr).Error(); err != nil {
		return false, err
	}

	if int64s.Contains(tr.Archived, int64(mailID)) {
		log.Printf("mails %d is archived", mailID)
		return false, nil
	} else if int64s.Contains(tr.Deleted, int64(mailID)) {
		log.Printf("mails %d is deleted", mailID)
		return false, nil
	}

	return true, nil
}
