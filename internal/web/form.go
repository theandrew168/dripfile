package web

import (
	"strings"
	"unicode/utf8"
)

// Based on:
// Let's Go - Chapter 8.5 (Alex Edwards)
type Form struct {
	Errors map[string]string
}

func (f *Form) Valid() bool {
	return len(f.Errors) == 0
}

func (f *Form) AddError(key, message string) {
	if f.Errors == nil {
		f.Errors = make(map[string]string)
	}

	f.Errors[key] = message
}

func (f *Form) CheckField(ok bool, key, message string) {
	if !ok {
		f.AddError(key, message)
	}
}

func NotBlank(value string) bool {
	return strings.TrimSpace(value) != ""
}

func MaxCharacters(value string, n int) bool {
	return utf8.RuneCountInString(value) <= n
}
