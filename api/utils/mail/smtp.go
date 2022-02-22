package mail

import (
	"crypto/tls"
	"fmt"
	"net"
	"net/mail"
	"net/smtp"
	"os"
)

// SendMail – отправка почты по протоколу SSL/TLS
func SendMail(toSender string, toAddress string, subj string, body string) error {
	to := mail.Address{toSender, toAddress}
	from := mail.Address{os.Getenv("APP_NAME"), os.Getenv("MAIL_USER")}

	// Настройка заголовков письма
	headers := make(map[string]string)
	headers["From"] = from.String()
	headers["To"] = to.String()
	headers["Subject"] = subj

	// Настройка сообщения
	message := ""
	for k, v := range headers {
		message += fmt.Sprintf("%s: %s\r\n", k, v)
	}
	message += "\r\n" + body

	// Соединение с SMTP-сервером
	servername := os.Getenv("MAIL_HOST")
	host, _, _ := net.SplitHostPort(servername)
	auth := smtp.PlainAuth("", os.Getenv("MAIL_USER"), os.Getenv("MAIL_PASS"), host)

	// Конфигурация TLS
	tlsconfig := &tls.Config{
		InsecureSkipVerify: true,
		ServerName:         host,
	}

	var c smtp.Client

	if os.Getenv("MAIL_TYPE") == "SSL" {
		conn, err := tls.Dial("tcp", servername, tlsconfig)
		if err != nil {
			return err
		}

		c, err := smtp.NewClient(conn, host)
		if err != nil {
			return err
		}

		fmt.Sprintf("Соединение установлено с %s", c.Text)
	}

	if os.Getenv("MAIL_TYPE") == "StartTLS" {
		c, err := smtp.Dial(servername)
		if err != nil {
			return err
		}
		c.StartTLS(tlsconfig)
	}

	// Авторизация
	if err := c.Auth(auth); err != nil {
		return err
	}

	// Настройка отправителя и получателя
	if err := c.Mail(from.Address); err != nil {
		return err
	}

	if err := c.Rcpt(to.Address); err != nil {
		return err
	}

	// Настройка передаваемых данных
	w, err := c.Data()
	if err != nil {
		return err
	}

	_, err = w.Write([]byte(message))
	if err != nil {
		return err
	}

	err = w.Close()
	if err != nil {
		return err
	}

	c.Quit()

	return nil
}
