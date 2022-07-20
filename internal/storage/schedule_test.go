package storage_test

import (
	"github.com/theandrew168/dripfile/internal/model"
	"github.com/theandrew168/dripfile/internal/test"
)

func mockSchedule(project model.Project) model.Schedule {
	schedule := model.NewSchedule(
		test.RandomString(8),
		test.RandomString(8),
		project,
	)
	return schedule
}
