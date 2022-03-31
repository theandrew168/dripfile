package mail

import (
	"strings"

	"github.com/theandrew168/dripfile/pkg/log"
)

type logMailer struct {
	logger log.Logger
}

func NewLogMailer(logger log.Logger) Mailer {
	m := logMailer{
		logger: logger,
	}
	return &m
}

func (m *logMailer) SendEmail(fromName, fromEmail, toName, toEmail, subject, body string) error {
	m.logger.Info("--- LogMailer Start ---\n")
	m.logger.Info("SendEmail:\n")
	m.logger.Info("From: %s (%s)\n", fromName, fromEmail)
	m.logger.Info("To:   %s (%s)\n", toName, toEmail)
	m.logger.Info("%s\n", subject)
	m.logger.Info("  %s\n", strings.Replace(body, "\n", "\n  ", -1))
	m.logger.Info("--- LogMailer Finish ---\n")
	return nil
}
