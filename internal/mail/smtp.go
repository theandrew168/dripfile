package mail

type smtpMailer struct {
}

func NewSMTPMailer(uri string) (Mailer, error) {
	// TODO: parse URI, pull out deets
	return nil, nil
}

// TODO: implement SMTP Mailer
func (m *smtpMailer) SendEmail(fromName, fromEmail, toName, toEmail, subject, body string) error {
	return nil
}
