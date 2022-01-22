package log

import (
	"log"
	"os"
)

type Logger interface {
	Info(format string, args ...interface{})
	Error(err error)
}

type logger struct {
	logger *log.Logger
}

func NewLogger() Logger {
	l := logger{
		logger: log.New(os.Stdout, "", 0),
	}
	return &l
}

func (l *logger) Info(format string, args ...interface{}) {
	l.logger.Printf(format, args...)
}

func (l *logger) Error(err error) {
	l.logger.Printf("error: %s\n", err)
}
