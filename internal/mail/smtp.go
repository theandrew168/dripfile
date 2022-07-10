package mail

type smtpMailer struct {
}

func NewSMTPMailer(uri string) (Mailer, error) {
	return nil, nil
}

func (m *smtpMailer) SendEmail(fromName, fromEmail, toName, toEmail, subject, body string) error {
	return nil
}
