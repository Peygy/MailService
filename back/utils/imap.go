package utils

import (
	"backend/internal/model"
	"bytes"
	"crypto/tls"
	"io"
	"log"
	"mime"
	"mime/multipart"
	"net/mail"
	"strings"

	"github.com/emersion/go-imap"
	"github.com/emersion/go-imap/client"
	"golang.org/x/net/html"
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
	seqSet.AddRange(1, mbox.Messages)

	section := imap.BodySectionName{}
	items := []imap.FetchItem{section.FetchItem()}

	messages := make(chan *imap.Message, 10)
	done := make(chan error, 1)
	go func() {
		done <- c.Fetch(seqSet, items, messages)
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
			log.Println(subject)

			// // Обработка получателей
			// receivers := parseAddressList(to)
			// if cc != "" {
			// 	receivers = append(receivers, parseAddressList(cc)...)
			// }

			body, err := extractEmailBody(reader)
			if err != nil {
				log.Println("Failed to extract mail body:", err)
				continue
			}

			mailRecord := model.Mail{
				Sender:  from,
				Subject: subject,
				Body:    string(body),
				IsRead:  false,
			}

			to = strings.TrimSpace(strings.Split(to, " ")[0])
			mailRecord.Receivers.Set(to)
			db.Create(&mailRecord)

			// Помечаем письмо для удаления
			delSeqSet := new(imap.SeqSet)
			delSeqSet.AddNum(msg.SeqNum)
			if err := c.Store(delSeqSet, imap.FormatFlagsOp(imap.AddFlags, true), []interface{}{imap.DeletedFlag}, nil); err != nil {
				log.Println("Failed to mark message for deletion:", err)
				continue
			}
		}
	}

	// Окончательно удаляем письма, помеченные флагом `\Deleted`
	if err := c.Expunge(nil); err != nil {
		log.Println("Failed to expunge messages:", err)
	}

	return <-done
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

// Функция для извлечения текстового или HTML тела письма
func extractEmailBody(msg *mail.Message) (string, error) {
	mediaType, params, err := mime.ParseMediaType(msg.Header.Get("Content-Type"))
	if err != nil {
		return "", err
	}

	// Если письмо не multipart, просто читаем тело
	if !strings.HasPrefix(mediaType, "multipart/") {
		body, err := io.ReadAll(msg.Body)
		return string(body), err
	}

	// Разбираем multipart
	mr := multipart.NewReader(msg.Body, params["boundary"])
	var plainTextBody, htmlBody string

	for {
		part, err := mr.NextPart()
		if err != nil {
			break
		}

		partMediaType, _, err := mime.ParseMediaType(part.Header.Get("Content-Type"))
		if err != nil {
			continue
		}

		partBody, err := io.ReadAll(part)
		if err != nil {
			continue
		}

		if strings.HasPrefix(partMediaType, "text/plain") {
			plainTextBody = string(partBody)
		} else if strings.HasPrefix(partMediaType, "text/html") {
			htmlBody = string(partBody)
		}
	}

	// Если есть HTML, то используем его, иначе берем plain text
	if htmlBody != "" {
		return stripHTML(htmlBody), nil
	}
	return plainTextBody, nil
}

// Функция для удаления HTML тегов
func stripHTML(htmlStr string) string {
	doc, err := html.Parse(strings.NewReader(htmlStr))
	if err != nil {
		return htmlStr
	}
	var buf bytes.Buffer
	var f func(*html.Node)
	f = func(n *html.Node) {
		if n.Type == html.TextNode {
			buf.WriteString(n.Data)
		}
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			f(c)
		}
	}
	f(doc)
	return buf.String()
}
