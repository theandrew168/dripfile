package postmark

import (
	"log"
	"strings"
)

type mockImpl struct {
	infoLog *log.Logger
}

func NewMock(infoLog *log.Logger) Interface {
	i := mockImpl{
		infoLog: infoLog,
	}
	return &i
}

func (i *mockImpl) SendEmail(fromName, fromEmail, toName, toEmail, subject, body string) error {
	i.infoLog.Printf("postmark.SendEmail:\n")
	i.infoLog.Printf("From: %s (%s)\n", fromName, fromEmail)
	i.infoLog.Printf("To:   %s (%s)\n", toName, toEmail)
	i.infoLog.Printf("%s\n", subject)
	i.infoLog.Printf("  %s\n", strings.Replace(body, "\n", "\n  ", -1))
	return nil
}
