package storage_test

import (
	"github.com/theandrew168/dripfile/src/core"
	"github.com/theandrew168/dripfile/src/random"
)

func mockLocation(project core.Project) core.Location {
	location := core.NewLocation(
		random.String(8),
		random.String(8),
		random.Bytes(8),
		project,
	)
	return location
}
