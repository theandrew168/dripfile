package storage_test

import (
	"github.com/theandrew168/dripfile/internal/core"
	"github.com/theandrew168/dripfile/internal/test"
)

func mockSchedule(project core.Project) core.Schedule {
	schedule := core.NewSchedule(
		test.RandomString(8),
		test.RandomString(8),
		project,
	)
	return schedule
}
