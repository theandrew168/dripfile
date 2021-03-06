package mail

import (
	"github.com/theandrew168/dripfile/internal/jsonlog"
)

type mockMailer struct {
	logger *jsonlog.Logger
}

func NewMockMailer(logger *jsonlog.Logger) (Mailer, error) {
	m := mockMailer{
		logger: logger,
	}
	return &m, nil
}

func (m *mockMailer) SendEmail(fromName, fromEmail, toName, toEmail, subject, body string) error {
	m.logger.Info("send email", map[string]string{
		"from_name":  fromName,
		"from_email": fromEmail,
		"to_name":    toName,
		"to_email":   toEmail,
		"subject":    subject,
		"body":       body,
	})
	return nil
}
