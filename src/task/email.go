package task

import (
	"context"
	"encoding/json"

	"github.com/hibiken/asynq"

	"github.com/theandrew168/dripfile/src/mail"
)

const (
	TypeEmailSend = "email:send"
)

type EmailSendPayload struct {
	FromName  string `json:"from_name"`
	FromEmail string `json:"from_email"`
	ToName    string `json:"to_name"`
	ToEmail   string `json:"to_email"`
	Subject   string `json:"subject"`
	Body      string `json:"body"`
}

func NewEmailSendTask(fromName, fromEmail, toName, toEmail, subject, body string) (*asynq.Task, error) {
	payload := EmailSendPayload{
		FromName:  fromName,
		FromEmail: fromEmail,
		ToName:    toName,
		ToEmail:   toEmail,
		Subject:   subject,
		Body:      body,
	}

	js, err := json.Marshal(payload)
	if err != nil {
		return nil, err
	}

	return asynq.NewTask(TypeEmailSend, js), nil
}

type EmailSendWorker struct {
	mailer mail.Mailer
}

func NewEmailSendWorker(mailer mail.Mailer) *EmailSendWorker {
	w := EmailSendWorker{
		mailer: mailer,
	}
	return &w
}

func (w *EmailSendWorker) ProcessTask(ctx context.Context, t *asynq.Task) error {
	var payload EmailSendPayload
	err := json.Unmarshal(t.Payload(), &payload)
	if err != nil {
		return err
	}

	err = w.mailer.SendEmail(
		payload.FromName,
		payload.FromEmail,
		payload.ToName,
		payload.ToEmail,
		payload.Subject,
		payload.Body,
	)
	if err != nil {
		return err
	}

	return nil
}
