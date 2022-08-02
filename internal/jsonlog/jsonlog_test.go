package jsonlog_test

import (
	"bytes"
	"errors"
	"strings"
	"testing"

	"github.com/theandrew168/dripfile/internal/jsonlog"
	"github.com/theandrew168/dripfile/internal/test"
)

func TestInfo(t *testing.T) {
	t.Parallel()

	tests := []struct {
		message    string
		properties map[string]string
		want       string
	}{
		{"", nil, `{"level":"INFO","message":""}`},
		{"test", nil, `{"level":"INFO","message":"test"}`},
		{
			"test",
			map[string]string{"foo": "bar"},
			`{"level":"INFO","message":"test","properties":{"foo":"bar"}}`,
		},
	}

	for _, tt := range tests {
		var b bytes.Buffer
		logger := jsonlog.New(&b)
		logger.Info(tt.message, tt.properties)

		got := strings.TrimSpace(b.String())
		test.AssertEqual(t, got, tt.want)
	}
}

func TestInfof(t *testing.T) {
	t.Parallel()

	tests := []struct {
		format string
		args   []any
		want   string
	}{
		{"", nil, `{"level":"INFO","message":""}`},
		{"test", nil, `{"level":"INFO","message":"test"}`},
		{"%s", []any{"test"}, `{"level":"INFO","message":"test"}`},
		{"%s %d", []any{"test", 42}, `{"level":"INFO","message":"test 42"}`},
	}

	for _, tt := range tests {
		var b bytes.Buffer
		logger := jsonlog.New(&b)
		logger.Infof(tt.format, tt.args...)

		got := strings.TrimSpace(b.String())
		test.AssertEqual(t, got, tt.want)
	}
}

func TestError(t *testing.T) {
	t.Parallel()

	tests := []struct {
		err        error
		properties map[string]string
		want       string
	}{
		{errors.New(""), nil, `{"level":"ERROR","message":""}`},
		{errors.New("test"), nil, `{"level":"ERROR","message":"test"}`},
		{
			errors.New("test"),
			map[string]string{"foo": "bar"},
			`{"level":"ERROR","message":"test","properties":{"foo":"bar"}}`,
		},
	}

	for _, tt := range tests {
		var b bytes.Buffer
		logger := jsonlog.New(&b)
		logger.Error(tt.err, tt.properties)

		got := strings.TrimSpace(b.String())
		test.AssertEqual(t, got, tt.want)
	}
}

func TestErrorf(t *testing.T) {
	t.Parallel()

	tests := []struct {
		format string
		args   []any
		want   string
	}{
		{"", nil, `{"level":"ERROR","message":""}`},
		{"test", nil, `{"level":"ERROR","message":"test"}`},
		{"%s", []any{"test"}, `{"level":"ERROR","message":"test"}`},
		{"%s %d", []any{"test", 42}, `{"level":"ERROR","message":"test 42"}`},
	}

	for _, tt := range tests {
		var b bytes.Buffer
		logger := jsonlog.New(&b)
		logger.Errorf(tt.format, tt.args...)

		got := strings.TrimSpace(b.String())
		test.AssertEqual(t, got, tt.want)
	}
}
