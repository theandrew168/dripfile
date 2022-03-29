package mail

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

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
