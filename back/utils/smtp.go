package utils

import (
	"backend/internal/model"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"log"
	"net/smtp"

	"github.com/jordan-wright/email"
)

func SendMailSMTP(mail model.Mail) error {
	var (
		smtpHost = GetEnv("SMTP_HOST", "")
		smtpUser = GetEnv("SMTP_USER", "")
		smtpPass = GetEnv("MAIL_PASS", "")
	)

	var receivers []string
	if err := json.Unmarshal(mail.Receivers.Bytes, &receivers); err != nil {
		return err
	}

	e := email.NewEmail()
	e.From = fmt.Sprintf("\"%s\" <%s>", mail.Sender, smtpUser)
	e.To = receivers
	e.Subject = fmt.Sprintf("Письмо из GoMail! %s", mail.Subject)
	e.Text = []byte(mail.Body)

	auth := smtp.PlainAuth("", smtpUser, smtpPass, smtpHost)

	if err := e.SendWithTLS(smtpHost+":465", auth, &tls.Config{
		InsecureSkipVerify: true,
		ServerName:         smtpHost,
	}); err != nil {
		log.Println("Failed to send email:", err)
		return err
	}

	log.Println("Email sent successfully")
	return nil
}
