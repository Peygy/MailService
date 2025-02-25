package utils

import (
	"fmt"
	"log"
	"net/smtp"
)

func SendMailSMTP(to string, subject string, body string) error {
	// Конфигурация для SMTP сервера
	smtpHost := "smtp.gmail.com"       // Укажите SMTP-сервер (например, для Gmail)
	smtpPort := "587"                  // Порт для TLS
	smtpUser := "your-email@gmail.com" // Ваш email
	smtpPass := "your-email-password"  // Ваш пароль для почтового аккаунта (или App Password)

	// Формирование сообщения
	from := smtpUser
	msg := fmt.Sprintf("From: %s\n", from)
	msg += fmt.Sprintf("To: %s\n", to)
	msg += fmt.Sprintf("Subject: %s\n\n", subject)
	msg += body

	// Авторизация
	auth := smtp.PlainAuth("", smtpUser, smtpPass, smtpHost)

	// Отправка письма
	err := smtp.SendMail(smtpHost+":"+smtpPort, auth, from, []string{to}, []byte(msg))
	if err != nil {
		log.Println("Failed to send email:", err)
		return err
	}

	log.Println("Email sent successfully")
	return nil
}
