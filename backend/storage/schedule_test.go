package storage_test

import (
	"github.com/theandrew168/dripfile/backend/core"
	"github.com/theandrew168/dripfile/backend/random"
)

func mockSchedule(project core.Project) core.Schedule {
	schedule := core.NewSchedule(
		random.String(8),
		random.String(8),
		project,
	)
	return schedule
}
