package storage_test

import (
	"github.com/theandrew168/dripfile/src/core"
	"github.com/theandrew168/dripfile/src/random"
)

func mockSchedule(project core.Project) core.Schedule {
	schedule := core.NewSchedule(
		random.String(8),
		random.String(8),
		project,
	)
	return schedule
}
