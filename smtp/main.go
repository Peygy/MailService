package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"net/smtp"
	"strings"
)

func main() {
	// Запуск SMTP-сервера на порту 25
	ln, err := net.Listen("tcp", ":25")
	if err != nil {
		log.Fatal("Ошибка при запуске сервера:", err)
	}
	defer ln.Close()

	log.Println("SMTP сервер запущен на порту 25...")

	for {
		conn, err := ln.Accept()
		if err != nil {
			log.Println("Ошибка при подключении:", err)
			continue
		}
		go handleConnection(conn)
	}
}

func handleConnection(conn net.Conn) {
	defer conn.Close()

	// Создаем буфер для чтения/записи
	reader := bufio.NewReader(conn)
	writer := bufio.NewWriter(conn)

	// Отправляем приветственное сообщение
	sendResponse(writer, "220 Simple SMTP Server Ready")

	// Чтение данных от клиента
	clientHello, _ := reader.ReadString('\n')
	fmt.Println("Received: ", clientHello)
	sendResponse(writer, "250 Hello")

	// Получаем команду MAIL FROM
	recipient, _ := reader.ReadString('\n')
	fmt.Println("Received MAIL FROM:", recipient)
	sendResponse(writer, "250 OK")

	// Получаем команду RCPT TO
	recipientTo, _ := reader.ReadString('\n')
	fmt.Println("Received RCPT TO:", recipientTo)
	sendResponse(writer, "250 OK")

	// Получаем DATA
	sendResponse(writer, "354 Start mail input; end with <CRLF>.<CRLF>")

	// Чтение самого письма
	var message []string
	for {
		line, _ := reader.ReadString('\n')
		if line == "\r\n" {
			break
		}
		message = append(message, line)
	}

	// Сохраняем письмо в лог (или отправляем)
	logMessage(message)

	// Отправляем подтверждение получения письма
	sendResponse(writer, "250 Message accepted")

	// Отправляем письмо через SMTP на указанный адрес
	sendMail(message)
}

func sendResponse(writer *bufio.Writer, response string) {
	writer.WriteString(response + "\r\n")
	writer.Flush()
}

func logMessage(message []string) {
	log.Printf("Received message:\n%s\n", strings.Join(message, "\n"))
}

// Функция для отправки письма на локальный SMTP-сервер
func sendMail(message []string) {
	// Параметры SMTP-сервера
	smtpHost := "localhost:25"  // Локальный SMTP-сервер на порту 25
	from := "sender@localhost"  // Отправитель
	to := "recipient@localhost" // Получатель

	// Формируем тело письма
	subject := "Subject: Test Email from Local Server"
	body := strings.Join(message, "\n")

	// Формируем письмо
	msg := []byte(subject + "\r\n" + body)

	// Отправляем письмо
	err := smtp.SendMail(smtpHost, nil, from, []string{to}, msg)
	if err != nil {
		log.Printf("Ошибка при отправке письма: %v", err)
		return
	}

	log.Println("Письмо успешно отправлено на", to)
}
