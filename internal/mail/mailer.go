package mail

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/theandrew168/dripfile/internal/log"
)

type Mailer interface {
	SendEmail(fromName, fromEmail, toName, toEmail, subject, body string) error
}

type postmarkMailer struct {
	postmarkAPIKey string
}

func NewPostmarkMailer(postmarkAPIKey string) Mailer {
	m := postmarkMailer{
		postmarkAPIKey: postmarkAPIKey,
	}
	return &m
}

func (m *postmarkMailer) SendEmail(fromName, fromEmail, toName, toEmail, subject, body string) error {
	message := struct {
		From     string `json:"From"`
		To       string `json:"To"`
		Subject  string `json:"Subject"`
		TextBody string `json:"TextBody"`
	}{
		From:     fromEmail,
		To:       toEmail,
		Subject:  subject,
		TextBody: body,
	}

	b, err := json.Marshal(message)
	if err != nil {
		return err
	}

	client := &http.Client{}
	endpoint := "https://api.postmarkapp.com/email"
	req, err := http.NewRequest("POST", endpoint, bytes.NewBuffer(b))
	if err != nil {
		return err
	}

	req.Header.Set("Accept", "application/json")
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Postmark-Server-Token", m.postmarkAPIKey)
	resp, err := client.Do(req)
	if err != nil {
		return err
	}

	// error if status isn't a 2xx
	if resp.Status[0] != '2' {
		return fmt.Errorf("failed to send email: %s", resp.Status)
	}

	return nil
}

type logMailer struct {
	logger log.Logger
}

func NewLogMailer(logger log.Logger) Mailer {
	m := logMailer{
		logger: logger,
	}
	return &m
}

func (m *logMailer) SendEmail(fromName, fromEmail, toName, toEmail, subject, body string) error {
	m.logger.Info("--- LogMailer Start ---\n")
	m.logger.Info("From: %s (%s)\n", fromName, fromEmail)
	m.logger.Info("To:   %s (%s)\n", toName, toEmail)
	m.logger.Info("%s\n", subject)
	m.logger.Info("  %s\n", strings.Replace(body, "\n", "\n  ", -1))
	m.logger.Info("--- LogMailer Finish ---\n")
	return nil
}
