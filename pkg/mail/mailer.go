package mail

import ()

type Mailer interface {
	SendEmail(fromName, fromEmail, toName, toEmail, subject, body string) error
}
