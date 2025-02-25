package utils

import (
	"backend/internal/model"
	"crypto/tls"
	"fmt"
	"io"
	"log"

	"github.com/emersion/go-imap"
	"github.com/emersion/go-imap/client"
	"github.com/emersion/go-message/mail"
	"gorm.io/gorm"
)

const (
	imapServer = "imap.yandex.ru:993"
	emailUser  = "isakov.29072004"
	emailPass  = "anuptulzrvloszth"
)

func ReadMailIMAP(db *gorm.DB) error {
	// Подключаемся к серверу IMAP
	c, err := client.DialTLS(imapServer, &tls.Config{})
	if err != nil {
		log.Println("Failed to connect to IMAP server:", err)
		return err
	}
	defer c.Logout()

	// Авторизация
	if err := c.Login(emailUser, emailPass); err != nil {
		log.Println("Failed to login to IMAP server:", err)
		return err
	}

	// Выбираем папку "INBOX"
	mbox, err := c.Select("INBOX", false)
	if err != nil {
		log.Println("Failed to select INBOX:", err)
		return err
	}
	log.Printf("Total messages in INBOX: %d\n", mbox.Messages)

	if mbox.Messages == 0 {
		log.Println("No messages found")
		return nil
	}

	// Указываем диапазон писем для чтения (последние 10)
	seqSet := new(imap.SeqSet)
	seqSet.AddRange(uint32(mbox.Messages-9), uint32(mbox.Messages))

	// Указываем какие части письма хотим получить
	section := imap.BodySectionName{}
	items := []imap.FetchItem{section.FetchItem()}

	// Запрашиваем письма
	messages := make(chan *imap.Message, 10)
	go func() {
		if err := c.Fetch(seqSet, items, messages); err != nil {
			log.Println("Failed to fetch messages:", err)
		}
	}()

	// Обрабатываем письма
	for msg := range messages {
		if msg == nil {
			log.Println("Server didn't return message")
			continue
		}
		for _, value := range msg.Body {
			r, err := mail.CreateReader(value)
			if err != nil {
				log.Println("Failed to create mail reader:", err)
				continue
			}
			header := r.Header
			from, _ := header.AddressList("From")
			to, _ := header.AddressList("To")
			subject, _ := header.Subject()
			fmt.Printf("From: %v\nTo: %v\nSubject: %s\n", from, to, subject)

			var bodyText string
			for {
				p, err := r.NextPart()
				if err != nil {
					break
				}
				b, _ := io.ReadAll(p.Body)
				bodyText += string(b) + "\n"
			}

			// Сохранение в базу данных
			mailRecord := model.Mail{
				Sender:   from[0].Address,
				Receiver: to[0].Address,
				Subject:  subject,
				Body:     bodyText,
				IsRead:   false,
			}
			db.Create(&mailRecord)
		}
	}
	return nil
}
