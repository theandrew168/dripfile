package postmark

type Interface interface {
	SendEmail(fromName, fromEmail, toName, toEmail, subject, body string) error
}
