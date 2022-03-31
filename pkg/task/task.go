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

	return New(KindEmail, string(b)), nil
}

type PruneSessionsInfo struct{}

func PruneSessions() (Task, error) {
	info := PruneSessionsInfo{}

	b, err := json.Marshal(info)
	if err != nil {
		return Task{}, err
	}

	return New(KindSession, string(b)), nil
}

type DoTransferInfo struct {
	ID string `json:"id"`
}

func DoTransfer(id string) (Task, error) {
	info := DoTransferInfo{
		ID: id,
	}

	b, err := json.Marshal(info)
	if err != nil {
		return Task{}, err
	}

	return New(KindTransfer, string(b)), nil
}
