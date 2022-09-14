package storage_test

import (
	"testing"
)

// used for record cleanup
type DeleterFunc func(t *testing.T)
