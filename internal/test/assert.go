package test

import (
	"strings"
	"testing"
)

func AssertEqual[T comparable](t *testing.T, got, want T) {
	t.Helper()

	if got != want {
		t.Fatalf("got %v; want %v", got, want)
	}
}

func AssertStringContains(t *testing.T, got, want string) {
	t.Helper()

	if !strings.Contains(got, want) {
		t.Fatalf("got %q; want to contain: %q", got, want)
	}
}

func AssertNilError(t *testing.T, got error) {
	t.Helper()

	if got != nil {
		t.Fatalf("got: %v; want: nil", got)
	}
}

func AssertErrorContains(t *testing.T, got error, want string) {
	t.Helper()

	if got == nil {
		t.Fatalf("got: nil; want: error to contain: %q", want)
	}

	if !strings.Contains(got.Error(), want) {
		t.Fatalf("got %q; want to contain: %q", got.Error(), want)
	}
}
