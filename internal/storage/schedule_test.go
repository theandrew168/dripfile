package storage_test

import (
	"github.com/theandrew168/dripfile/internal/model"
	"github.com/theandrew168/dripfile/internal/test"
)

func mockSchedule() model.Schedule {
	schedule := model.NewSchedule(
		test.RandomString(8),
		test.RandomString(8),
	)
	return schedule
}
