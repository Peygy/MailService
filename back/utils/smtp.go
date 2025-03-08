package utils

import (
	"backend/internal/model"
	"crypto/tls"
	"errors"
	"fmt"
	"log"
	"net/smtp"
	"strings"

	"github.com/jordan-wright/email"
)

const (
	domain = "gomail.kurs"
)

func SendMailSMTP(mail model.Mail, recs []string) error {
	var (
		smtpHost = GetEnv("SMTP_HOST", "")
		smtpUser = GetEnv("SMTP_USER", "")
		smtpPass = GetEnv("MAIL_PASS", "")
	)

	var filtered []string
	for _, rec := range recs {
		if !strings.Contains(rec, domain) {
			filtered = append(filtered, rec)
		}
	}

	if len(filtered) <= 0 {
		err := errors.New("receivers are empty")
		log.Println("Receivers are empty:", err)
		return err
	}

	e := email.NewEmail()
	e.From = fmt.Sprintf("\"%s\" <%s>", mail.Sender, smtpUser)
	e.To = filtered
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
