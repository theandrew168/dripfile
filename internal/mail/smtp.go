package mail

import (
//	"net/smtp"
//	"net/url"
)

type smtpMailer struct {
	hostname string
	port     string
	username string
	password string
}

func NewSMTPMailer(smtpURL string) (Mailer, error) {
	// TODO: parse URI, pull out deets
	return nil, nil
}

// TODO: implement SMTP Mailer
func (m *smtpMailer) SendEmail(fromName, fromEmail, toName, toEmail, subject, body string) error {
	return nil
}
