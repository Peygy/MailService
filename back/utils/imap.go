package utils

import (
	"backend/internal/model"
	"crypto/tls"
	"io"
	"log"
	"net/mail"
	"strings"

	"github.com/emersion/go-imap"
	"github.com/emersion/go-imap/client"
	"gorm.io/gorm"
)

func ReadMailIMAP(db *gorm.DB) error {
	var (
		imapHost = GetEnv("IMAP_HOST", "")
		imapUser = GetEnv("IMAP_USER", "")
		imapPass = GetEnv("MAIL_PASS", "")
	)

	c, err := client.DialTLS(imapHost, &tls.Config{})
	if err != nil {
		log.Println("Failed to connect to IMAP server:", err)
		return err
	}
	defer c.Logout()

	if err := c.Login(imapUser, imapPass); err != nil {
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
			// cc := header.Get("Cc")
			subject := header.Get("Subject")

			log.Println(from)
			log.Println(to)
			// log.Println(cc)
			log.Println(subject)

			// // Обработка получателей
			// receivers := parseAddressList(to)
			// if cc != "" {
			// 	receivers = append(receivers, parseAddressList(cc)...)
			// }

			body, err := io.ReadAll(reader.Body)
			if err != nil {
				log.Println("Failed to read mail body:", err)
				continue
			}

			mailRecord := model.Mail{
				Sender:    from,
				Subject:   subject,
				Body:      string(body),
				IsRead:    false,
			}

			to = (strings.Split(to, " ")[0])
			to = strings.Trim(to, "\n")

			mailRecord.Receivers.Set(to)
			db.Create(&mailRecord)
		}
	}

	return nil
}

// func parseAddressList(addressList string) []string {
// 	var addresses []string
// 	parsedAddresses, err := mail.ParseAddressList(addressList)
// 	if err != nil {
// 		log.Println("Failed to parse address list:", err)
// 		return addresses
// 	}
// 	for _, addr := range parsedAddresses {
// 		addresses = append(addresses, addr.Address)
// 	}
// 	return addresses
// }
