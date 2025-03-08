package utils

import (
	"backend/internal/model"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"log"
	"net/smtp"
	"strings"

	"github.com/jordan-wright/email"
)

const (
	domain = "gomail.kurs"
)

type mailData struct {
	Receivers []string `json:"receivers"`
	Subject   string   `json:"subject"`
	Body      string   `json:"body"`
}

func SendMailSMTP(mail model.Mail, data mailData) error {
	var (
		smtpHost = GetEnv("SMTP_HOST", "")
		smtpUser = GetEnv("SMTP_USER", "")
		smtpPass = GetEnv("MAIL_PASS", "")
	)

	var filteredReceivers []string
	for _, rec := range data.Receivers {
		if !strings.Contains(rec, domain) {
			filteredReceivers = append(filteredReceivers, rec)
		}
	}
	data.Receivers = filteredReceivers
	if len(data.Receivers) <= 0 {
		return nil
	}

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
