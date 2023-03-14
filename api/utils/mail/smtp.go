// Package mail - пакет для отправки писем
package mail

import (
	"crypto/tls"
	"fmt"
	"net"
	"net/mail"
	"net/smtp"
	"os"
	"time"

	"github.com/doka-guide/api/api/utils/randomize"
)

// SendMail – отправка письма по SSL/TLS соединению
func SendMail(toSender string, toAddress string, subj string, textBody string, htmlBody string, isBulk bool) error {
	to := mail.Address{Name: toSender, Address: toAddress}
	from := mail.Address{
		Name:    os.Getenv("MAIL_SENDER"),
		Address: os.Getenv("MAIL_USER"),
	}

	// Формирование уникального разделителя
	boundary := randomize.GetRandomString(64)

	// Формирование заголовков
	headers := make(map[string]string)
	headers["MIME-Version"] = "1.0"
	headers["Date"] = time.Now().Format(time.RFC1123Z)
	headers["From"] = from.String()
	headers["To"] = to.String()
	headers["Subject"] = subj
	if isBulk {
		headers["Precedence"] = "bulk"
		headers["Reply-To"] = os.Getenv("MAIL_SENDER")
		headers["List-Unsubscribe"] = "<mailto:" + os.Getenv("MAIL_USER") + ">, <" + os.Getenv("MAIL_URL") + ">"
	}
	headers["Content-Type"] = "multipart/alternative; boundary=\"" + boundary + "\""
	headers["X-Sender"] = os.Getenv("MAIL_SENDER")
	headers["User-Agent"] = "Doka API"

	// Формирование заголовков письма
	message := ""
	for k, v := range headers {
		message += fmt.Sprintf("%s: %s\r\n", k, v)
		delete(headers, k)
	}

	// Формирование текстовой части сообщения
	headers["Content-Type"] = "text/plain; charset=utf-8"
	headers["Content-Transfer-Encoding"] = "quoted-printable"
	headers["Content-Disposition"] = "inline"

	message += "\r\n" + "--" + boundary + "\r\n"
	for k, v := range headers {
		message += fmt.Sprintf("%s: %s\r\n", k, v)
		delete(headers, k)
	}
	message += "\r\n" + textBody + "\r\n"

	// Формирование HTML части сообщения
	headers["Content-Type"] = "text/html; charset=\"utf-8\""
	headers["Content-Transfer-Encoding"] = "quoted-printable"
	headers["Content-Disposition"] = "inline"

	message += "\r\n" + "--" + boundary + "\r\n"
	for k, v := range headers {
		message += fmt.Sprintf("%s: %s\r\n", k, v)
		delete(headers, k)
	}
	message += "\r\n" + htmlBody + "\r\n"

	// Последняя строчка письма
	message += "\r\n" + "--" + boundary + "--" + "\r\n"

	// Соединение с SMTP-сервером
	serverName := os.Getenv("MAIL_HOST")
	host, _, _ := net.SplitHostPort(serverName)
	auth := smtp.PlainAuth(
		"Hi",
		os.Getenv("MAIL_USER"),
		os.Getenv("MAIL_PASS"),
		host,
	)

	// Конфигурация TLS
	tlsConfig := &tls.Config{
		InsecureSkipVerify: true,
		ServerName:         host,
	}

	// Соединение с сервером (no StartTLS)
	conn, err := tls.Dial("tcp", serverName, tlsConfig)
	if err != nil {
		return err
	}

	c, err := smtp.NewClient(conn, host)
	if err != nil {
		return err
	}

	// Авторизация
	if err = c.Auth(auth); err != nil {
		return err
	}

	// Формирование отправителя и получателя
	if err = c.Mail(from.Address); err != nil {
		return err
	}

	if err = c.Rcpt(to.Address); err != nil {
		return err
	}

	// Подготовка данных
	w, err := c.Data()
	if err != nil {
		return err
	}

	// Пересылка данных
	_, err = w.Write([]byte(message))
	if err != nil {
		return err
	}

	// Закрытие соединения
	err = w.Close()
	if err != nil {
		return err
	}

	// Выход из аккаунта
	c.Quit()

	return nil
}
