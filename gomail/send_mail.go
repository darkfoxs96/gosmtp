package gomail

import (
	"bytes"
	"errors"
	"fmt"
	"mime/quotedprintable"
	"net/smtp"
	"strings"
)

/**
	Modified from https://gist.github.com/jpillora/cb46d183eca0710d909a
	Thank you very much.
**/

// SMTPServer stores host, port for SMTP map
type SMTPServer struct {
	Host string
	Port string
}

// SMTP Servers
var SMTP = map[string]*SMTPServer{
	"mail.ru":   &SMTPServer{Host: "smtp.mail.ru", Port: "465"},
	"yandex.ru": &SMTPServer{Host: "smtp.yandex.ru", Port: "465"},
	"gmail.com": &SMTPServer{Host: "smtp.gmail.com", Port: "465"},
}

// Sender client request
type Sender struct {
	User     string
	Password string
	SMTPServer
}

// NewSender create Sender
func NewSender(Username, Password, SMTPHost, SMTPPort string) (sender Sender, err error) {
	if SMTPHost == "" || SMTPPort == "" || Password == "" || Username == "" {
		err = errors.New("gomail: have empty field")
	}
	sender = Sender{Username, Password, SMTPServer{SMTPHost, SMTPPort}}
	return
}

// SendMail send mail to client
func (sender Sender) SendMail(Dest []string, Subject, bodyMessage string) (err error) {
	msg := "From: " + sender.User + "\n" +
		"To: " + strings.Join(Dest, ",") + "\n" +
		"Subject: " + Subject + "\n" + bodyMessage

	err = smtp.SendMail(sender.SMTPServer.Host+":"+sender.SMTPServer.Port,
		smtp.PlainAuth("", sender.User, sender.Password, sender.SMTPServer.Host),
		sender.User, Dest, []byte(msg))

	if err != nil {
		return
	}
	return
}

// WriteEmail to message
func (sender Sender) WriteEmail(dest []string, contentType, subject, bodyMessage string) string {
	header := make(map[string]string)
	header["From"] = sender.User

	receipient := ""

	for _, user := range dest {
		receipient = receipient + user
	}

	header["To"] = receipient
	header["Subject"] = subject
	header["MIME-Version"] = "1.0"
	header["Content-Type"] = fmt.Sprintf("%s; charset=\"utf-8\"", contentType)
	header["Content-Transfer-Encoding"] = "quoted-printable"
	header["Content-Disposition"] = "inline"

	message := ""

	for key, value := range header {
		message += fmt.Sprintf("%s: %s\r\n", key, value)
	}

	var encodedMessage bytes.Buffer

	finalMessage := quotedprintable.NewWriter(&encodedMessage)
	finalMessage.Write([]byte(bodyMessage))
	finalMessage.Close()

	message += "\r\n" + encodedMessage.String()

	return message
}

// WriteHTMLEmail to html email
func (sender *Sender) WriteHTMLEmail(dest []string, subject, bodyMessage string) string {
	return sender.WriteEmail(dest, "text/html", subject, bodyMessage)
}

// WritePlainEmail to plain email
func (sender *Sender) WritePlainEmail(dest []string, subject, bodyMessage string) string {
	return sender.WriteEmail(dest, "text/plain", subject, bodyMessage)
}
