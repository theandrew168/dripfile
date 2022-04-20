package postmark

import (
	"github.com/theandrew168/dripfile/pkg/jsonlog"
)

type mockImpl struct {
	logger *jsonlog.Logger
}

func NewMock(logger *jsonlog.Logger) Interface {
	i := mockImpl{
		logger: logger,
	}
	return &i
}

func (i *mockImpl) SendEmail(fromName, fromEmail, toName, toEmail, subject, body string) error {
	i.logger.PrintInfo("postmark email send", map[string]string{
		"from_name":  fromName,
		"from_email": fromEmail,
		"to_name":    toName,
		"to_email":   toEmail,
		"subject":    subject,
		"body":       body,
	})
	return nil
}
