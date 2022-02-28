package task

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
}

func New(kind, info string) Task {
	task := Task{
		Kind:   kind,
		Info:   info,
		Status: StatusNew,
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

func NewEmailInfo(fromName, fromEmail, toName, toEmail, subject, body string) EmailInfo {
	info := EmailInfo{
		FromName:  fromName,
		FromEmail: fromEmail,
		ToName:    toName,
		ToEmail:   toEmail,
		Subject:   subject,
		Body:      body,
	}
	return info
}

type SessionInfo struct{}

func NewSessionInfo() SessionInfo {
	info := SessionInfo{}
	return info
}

type TransferInfo struct {
	ID string `json:"id"`
}

func NewTransferInfo(id string) TransferInfo {
	info := TransferInfo{
		ID: id,
	}
	return info
}
