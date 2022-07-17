package storage_test

import (
	"github.com/theandrew168/dripfile/internal/core"
	"github.com/theandrew168/dripfile/internal/test"
)

func mockAccount(project core.Project) core.Account {
	account := core.NewAccount(
		test.RandomString(8),
		test.RandomString(8),
		test.RandomString(8),
		project,
	)
	return account
}
