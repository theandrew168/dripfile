package jsonlog

import (
	"encoding/json"
	"fmt"
	"io"
	"sync"
)

const (
	LevelInfo  = "INFO"
	LevelError = "ERROR"
)

type Logger struct {
	mu  sync.Mutex
	out io.Writer
}

func New(out io.Writer) *Logger {
	l := Logger{
		out: out,
	}
	return &l
}

func (l *Logger) Info(message string, properties map[string]string) {
	l.print(LevelInfo, message, properties)
}

func (l *Logger) Infof(format string, args ...any) {
	message := fmt.Sprintf(format, args...)
	l.print(LevelInfo, message, nil)
}

func (l *Logger) Error(err error, properties map[string]string) {
	l.print(LevelError, err.Error(), properties)
}

func (l *Logger) Errorf(format string, args ...any) {
	message := fmt.Sprintf(format, args...)
	l.print(LevelError, message, nil)
}

func (l *Logger) print(level, message string, properties map[string]string) (int, error) {
	// use an anon struct to represent the log entry
	entry := struct {
		Level      string            `json:"level"`
		Message    string            `json:"message"`
		Properties map[string]string `json:"properties,omitempty"`
	}{
		Level:      level,
		Message:    message,
		Properties: properties,
	}

	var line []byte
	line, err := json.Marshal(entry)
	if err != nil {
		line = []byte(LevelError + ": unable to marshal log message: " + err.Error())
	}

	l.mu.Lock()
	defer l.mu.Unlock()

	return l.out.Write(append(line, '\n'))
}
