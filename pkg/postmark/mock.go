package postmark

import (
	"strings"

	"github.com/theandrew168/dripfile/pkg/log"
)

type mockImpl struct {
	logger log.Logger
}

func NewMock(logger log.Logger) Interface {
	i := mockImpl{
		logger: logger,
	}
	return &i
}

func (i *mockImpl) SendEmail(fromName, fromEmail, toName, toEmail, subject, body string) error {
	i.logger.Info("postmark.SendEmail:\n")
	i.logger.Info("From: %s (%s)\n", fromName, fromEmail)
	i.logger.Info("To:   %s (%s)\n", toName, toEmail)
	i.logger.Info("%s\n", subject)
	i.logger.Info("  %s\n", strings.Replace(body, "\n", "\n  ", -1))
	return nil
}
