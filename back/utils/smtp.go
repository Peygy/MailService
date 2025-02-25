package utils

import (
	"backend/internal/model"
	"fmt"
	"log"
	"net/smtp"
)

func SendMailSMTP(mail model.Mail) error {
	// Конфигурация для SMTP сервера
	smtpHost := "smtp.yandex.ru"
	smtpPort := "465"
	smtpUser := "isakov.29072004"
	smtpPass := "anuptulzrvloszth"

	msg := fmt.Sprintf("From: %s\n", mail.Sender)
	msg += fmt.Sprintf("To: %s\n", mail.Receiver)
	msg += fmt.Sprintf("Subject: %s\n\n", mail.Subject)
	msg += mail.Body

	// Авторизация
	auth := smtp.PlainAuth("", smtpUser, smtpPass, smtpHost)

	// Отправка письма
	err := smtp.SendMail(smtpHost+":"+smtpPort, auth, mail.Sender, []string{mail.Receiver}, []byte(msg))
	if err != nil {
		log.Println("Failed to send email:", err)
		return err
	}

	log.Println("Email sent successfully")
	return nil
}
