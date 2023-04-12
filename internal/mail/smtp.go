package mail

import (
	"fmt"
	"net/smtp"
	"net/url"
)

type smtpMailer struct {
	username string
	password string
	host     string
	port     string
}

func NewSMTPMailer(smtpURI string) (Mailer, error) {
	u, err := url.Parse(smtpURI)
	if err != nil {
		return nil, err
	}

	username := u.User.Username()
	password, _ := u.User.Password()
	m := smtpMailer{
		username: username,
		password: password,
		host:     u.Hostname(),
		port:     u.Port(),
	}
	return &m, nil
}

func (m *smtpMailer) SendEmail(fromName, fromEmail, toName, toEmail, subject, body string) error {
	var from string
	if fromName != "" {
		from = fmt.Sprintf("%s <%s>", fromName, fromEmail)
	} else {
		from = fromEmail
	}

	var to string
	if toName != "" {
		to = fmt.Sprintf("%s <%s>", toName, toEmail)
	} else {
		to = toEmail
	}

	headers := map[string]string{
		"From":    from,
		"To":      to,
		"Subject": subject,
	}

	var header string
	for k, v := range headers {
		header += fmt.Sprintf("%s: %s\r\n", k, v)
	}

	// The msg parameter should be an RFC 822-style email with headers first,
	// a blank line, and then the message body. The lines of msg should be
	// CRLF terminated.
	msg := fmt.Sprintf("%s\r\n%s\r\n", header, body)

	addr := fmt.Sprintf("%s:%s", m.host, m.port)
	auth := smtp.PlainAuth("", m.username, m.password, m.host)
	err := smtp.SendMail(addr, auth, fromEmail, []string{toEmail}, []byte(msg))
	if err != nil {
		return err
	}

	return nil
}
