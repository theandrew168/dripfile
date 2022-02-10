package mail

type Info struct {
	FromName  string
	FromEmail string
	ToName    string
	ToEmail   string
	Body      string
}

type Mailer interface {
	SendMail(info Info) error
}
