package storage_test

import (
	"github.com/theandrew168/dripfile/internal/model"
	"github.com/theandrew168/dripfile/internal/test"
)

func mockLocation(project model.Project) model.Location {
	location := model.NewLocation(
		test.RandomString(8),
		test.RandomString(8),
		test.RandomBytes(8),
		project,
	)
	return location
}
