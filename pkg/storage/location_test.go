package storage_test

import (
	"github.com/theandrew168/dripfile/pkg/core"
	"github.com/theandrew168/dripfile/pkg/random"
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
