package utils

import (
	"backend/internal/model"
	"crypto/tls"
	"fmt"
	"io"
	"log"
	"net/mail"

	"github.com/emersion/go-imap"
	"github.com/emersion/go-imap/client"
	"gorm.io/gorm"
)

const (
	imapServer = "imap.yandex.ru:993"
	emailUser  = "isakov.29072004"
	emailPass  = "anuptulzrvloszth"
)

func ReadMailIMAP(db *gorm.DB) error {
	imapServer := imapServer
	emailUser := emailUser
	emailPass := emailPass
	if emailUser == "" || emailPass == "" {
		log.Println("IMAP credentials are missing")
		return fmt.Errorf("IMAP credentials are missing")
	}
	c, err := client.DialTLS(imapServer, &tls.Config{})
	if err != nil {
		log.Println("Failed to connect to IMAP server:", err)
		return err
	}
	defer c.Logout()
	if err := c.Login(emailUser, emailPass); err != nil {
		log.Println("Failed to login to IMAP server:", err)
		return err
	}
	mbox, err := c.Select("INBOX", false)
	if err != nil {
		log.Println("Failed to select INBOX:", err)
		return err
	}
	if mbox.Messages == 0 {
		log.Println("No messages found")
		return nil
	}

	seqSet := new(imap.SeqSet)
	seqSet.AddRange(mbox.Messages-9, mbox.Messages)

	section := imap.BodySectionName{}
	items := []imap.FetchItem{section.FetchItem()}

	messages := make(chan *imap.Message, 10)
	log.Println("IMAP reciever working...")
	go func() {
		if err := c.Fetch(seqSet, items, messages); err != nil {
			log.Println("Failed to fetch messages:", err)
		}
	}()

	for msg := range messages {
		if msg == nil {
			log.Println("Server didn't return message")
			continue
		}

		for _, value := range msg.Body {
			reader, err := mail.ReadMessage(value)
			if err != nil {
				log.Println("Failed to read mail message:", err)
				continue
			}

			header := reader.Header
			from := header.Get("From")
			to := header.Get("To")
			subject := header.Get("Subject")

			body, err := io.ReadAll(reader.Body)
			if err != nil {
				log.Println("Failed to read mail body:", err)
				continue
			}

			mailRecord := model.Mail{
				Sender:   from,
				Receiver: to,
				Subject:  subject,
				Body:     string(body),
				IsRead:   false,
			}
			db.Create(&mailRecord)
		}
	}

	return nil
}
