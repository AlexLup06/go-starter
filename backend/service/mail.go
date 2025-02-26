package service

import "net/smtp"

type MailService struct {
	apiKey string
}

func NewMailService(key string) *MailService {
	return &MailService{apiKey: key}
}

func (m *MailService) SendMail(from, to, subject, body string) error {
	// SMTP Server settings
	smtpServer := "smtp.sendgrid.net"
	smtpPort := "587"
	smtpAuth := smtp.PlainAuth("", "apikey", m.apiKey, smtpServer)

	// Email message
	msg := []byte(
		"From: " + from + "\r\n" +
			"To: " + to + "\r\n" +
			"Subject: " + subject + "\r\n" +
			"MIME-Version: 1.0\r\n" +
			"Content-Type: text/html; charset=\"utf-8\"\r\n" +
			"\r\n" + body)

	// Send email
	err := smtp.SendMail(smtpServer+":"+smtpPort, smtpAuth, from, []string{to}, msg)
	if err != nil {
		return err
	}
	return nil
}
