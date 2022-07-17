package validator

import (
	"fmt"
	"strings"
	"unicode/utf8"
)

// Based on:
// Let's Go - Chapter 8.5 (Alex Edwards)
type Validator struct {
	// general input error
	Error string

	// individual field errors
	FieldError map[string]string
}

func (v *Validator) Valid() bool {
	return len(v.FieldError) == 0
}

func (v *Validator) SetError(message string) {
	v.Error = message
}

func (v *Validator) SetFieldError(key, message string) {
	if v.FieldError == nil {
		v.FieldError = make(map[string]string)
	}

	v.FieldError[key] = message
}

func (v *Validator) CheckField(ok bool, key, message string) {
	if !ok {
		v.SetFieldError(key, message)
	}
}

func (v *Validator) CheckRequired(value, key string) {
	message := "This field is required"
	v.CheckField(Required(value), key, message)
}

func (v *Validator) CheckMaxCharacters(value string, n int, key string) {
	message := fmt.Sprintf("This field cannot me more than %d characters long", n)
	v.CheckField(MaxCharacters(value, n), key, message)
}

func Required(value string) bool {
	return strings.TrimSpace(value) != ""
}

func MaxCharacters(value string, n int) bool {
	return utf8.RuneCountInString(value) <= n
}
