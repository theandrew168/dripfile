package storage_test

import (
	"github.com/theandrew168/dripfile/internal/model"
	"github.com/theandrew168/dripfile/internal/test"
)

func mockAccount(project model.Project) model.Account {
	account := model.NewAccount(
		test.RandomString(8),
		test.RandomString(8),
		test.RandomString(8),
		project,
	)
	return account
}
