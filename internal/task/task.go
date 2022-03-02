package task

import (
	"encoding/json"
)

const (
	KindEmail    = "email"
	KindSession  = "session"
	KindTransfer = "transfer"
)

const (
	StatusNew     = "new"
	StatusRunning = "running"
	StatusSuccess = "success"
	StatusError   = "error"
)

type Task struct {
	// readonly (from database, after creation)
	ID string

	Kind   string
	Info   string
	Status string
	Error  string
}

func New(kind, info string) Task {
	task := Task{
		Kind:   kind,
		Info:   info,
		Status: StatusNew,
		Error:  "",
	}
	return task
}

type EmailInfo struct {
	FromName  string `json:"from_name"`
	FromEmail string `json:"from_email"`
	ToName    string `json:"to_name"`
	ToEmail   string `json:"to_email"`
	Subject   string `json:"subject"`
	Body      string `json:"body"`
}

func NewEmail(fromName, fromEmail, toName, toEmail, subject, body string) (Task, error) {
	info := EmailInfo{
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

	return New(KindEmail, string(b)), nil
}

type SessionInfo struct{}

func NewSession() (Task, error) {
	info := SessionInfo{}

	b, err := json.Marshal(info)
	if err != nil {
		return Task{}, err
	}

	return New(KindSession, string(b)), nil
}

type TransferInfo struct {
	ID string `json:"id"`
}

func NewTransfer(id string) (Task, error) {
	info := TransferInfo{
		ID: id,
	}

	b, err := json.Marshal(info)
	if err != nil {
		return Task{}, err
	}

	return New(KindTransfer, string(b)), nil
}
