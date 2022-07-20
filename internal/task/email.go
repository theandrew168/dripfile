package task

import (
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
