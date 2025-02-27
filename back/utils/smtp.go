package utils

import (
	"backend/internal/model"
	"crypto/tls"
	"fmt"
	"log"
	"net/smtp"

	"github.com/jordan-wright/email"
)

func SendMailSMTP(mail model.Mail) error {
	// Конфигурация для SMTP сервера
	smtpHost := "smtp.yandex.ru"
	smtpPort := "465"
	smtpUser := "isakov.29072004@yandex.ru"
	smtpPass := "anuptulzrvloszth"

	if smtpUser == "" || smtpPass == "" {
		log.Println("SMTP credentials are missing")
		return fmt.Errorf("SMTP credentials are missing")
	}

	e := email.NewEmail()
	e.From = fmt.Sprintf("%s <%s>", mail.Sender, smtpUser)
	e.To = []string{mail.Receiver}
	e.Subject = mail.Subject
	e.Text = []byte(mail.Body)

	auth := smtp.PlainAuth("", smtpUser, smtpPass, smtpHost)

	err := e.SendWithTLS(smtpHost+":"+smtpPort, auth, &tls.Config{
		InsecureSkipVerify: true,
		ServerName:         smtpHost,
	})
	if err != nil {
		log.Println("Failed to send email:", err)
		return err
	}

	log.Println("Email sent successfully")
	return nil
}
