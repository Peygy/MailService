package utils

import (
	"backend/internal/model"
	"crypto/tls"
	"errors"
	"fmt"
	"log"
	"net/smtp"

	"github.com/jordan-wright/email"
)

func SendMailSMTP(mail model.Mail, recs []string) error {
	var (
		smtpHost = GetEnv("SMTP_HOST", "")
		smtpUser = GetEnv("SMTP_USER", "")
		smtpPass = GetEnv("MAIL_PASS", "")
	)

	if len(recs) <= 0 {
		err := errors.New("receivers are empty")
		log.Println("Receivers are empty:", err)
		return nil
	}

	e := email.NewEmail()
	e.From = fmt.Sprintf("\"%s\" <%s>", mail.Sender, smtpUser)
	e.To = recs
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
