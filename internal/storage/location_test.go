package storage_test

import (
	"github.com/theandrew168/dripfile/internal/core"
	"github.com/theandrew168/dripfile/internal/test"
)

func mockLocation(project core.Project) core.Location {
	location := core.NewLocation(
		test.RandomString(8),
		test.RandomString(8),
		test.RandomBytes(8),
		project,
	)
	return location
}
