package task

import (
	"context"
	"encoding/json"
)

const (
	KindEmailSend = "email:send"
)

type EmailSendInfo struct {
	FromName  string `json:"from_name"`
	FromEmail string `json:"from_email"`
	ToName    string `json:"to_name"`
	ToEmail   string `json:"to_email"`
	Subject   string `json:"subject"`
	Body      string `json:"body"`
}

func NewEmailSendTask(fromName, fromEmail, toName, toEmail, subject, body string) Task {
	info := EmailSendInfo{
		FromName:  fromName,
		FromEmail: fromEmail,
		ToName:    toName,
		ToEmail:   toEmail,
		Subject:   subject,
		Body:      body,
	}

	js, err := json.Marshal(info)
	if err != nil {
		panic(err)
	}

	return NewTask(KindEmailSend, string(js))
}

func (w *Worker) HandleEmailSend(ctx context.Context, t Task) error {
	var info EmailSendInfo
	err := json.Unmarshal([]byte(t.Info), &info)
	if err != nil {
		return err
	}

	err = w.mailer.SendEmail(
		info.FromName,
		info.FromEmail,
		info.ToName,
		info.ToEmail,
		info.Subject,
		info.Body,
	)
	if err != nil {
		return err
	}

	return nil
}
