package task

import (
	"encoding/json"
)

const KindSendEmail = "send_email"

type SendEmailInfo struct {
	FromName  string `json:"from_name"`
	FromEmail string `json:"from_email"`
	ToName    string `json:"to_name"`
	ToEmail   string `json:"to_email"`
	Subject   string `json:"subject"`
	Body      string `json:"body"`
}

func SendEmail(fromName, fromEmail, toName, toEmail, subject, body string) (Task, error) {
	info := SendEmailInfo{
		FromName:  fromName,
		FromEmail: fromEmail,
		ToName:    toName,
		ToEmail:   toEmail,
		Subject:   subject,
		Body:      body,
	}

	b, err := json.Marshal(info)
	if err != nil {
		return Task{}, err
	}

	return New(KindSendEmail, string(b)), nil
}

func (w *Worker) SendEmail(task Task) error {
	var info SendEmailInfo
	err := json.Unmarshal([]byte(task.Info), &info)
	if err != nil {
		return err
	}

	err = w.postmark.SendEmail(
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
