package utils

import (
	"fmt"

	"gopkg.in/gomail.v2"
)

type EmailSender interface {
	SendEmail(subject string, content string, to []string) error
}

type GmailSender struct {
	name              string
	fromEmailAddress  string
	fromEmailPassword string
}

func NewGmailSender(name string, fromEmailAddress string, fromEmailPassword string) EmailSender {
	return &GmailSender{
		name:              name,
		fromEmailAddress:  fromEmailAddress,
		fromEmailPassword: fromEmailPassword,
	}
}

func (sender *GmailSender) SendEmail(subject string, content string, to []string) error {
	m := gomail.NewMessage()
	m.SetHeader("From", fmt.Sprintf("%s <%s>", sender.name, sender.fromEmailAddress))
	m.SetHeader("To", to...)
	m.SetHeader("Subject", subject)
	m.SetBody("text/html", content)

	// Gmail SMTP Server
	d := gomail.NewDialer("smtp.gmail.com", 587, sender.fromEmailAddress, sender.fromEmailPassword)

	return d.DialAndSend(m)
}