package mail

import (
	"golang.org/x/exp/slog"
)

type mockMailer struct {
	logger *slog.Logger
}

func NewMockMailer(logger *slog.Logger) (Mailer, error) {
	m := mockMailer{
		logger: logger,
	}
	return &m, nil
}

func (m *mockMailer) SendEmail(fromName, fromEmail, toName, toEmail, subject, body string) error {
	m.logger.Info("send email",
		slog.String("from_name", fromName),
		slog.String("from_email", fromEmail),
		slog.String("to_name", toName),
		slog.String("to_email", toEmail),
		slog.String("subject", subject),
		slog.String("body", body),
	)
	return nil
}
