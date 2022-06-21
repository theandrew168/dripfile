package storage_test

import (
	"github.com/theandrew168/dripfile/internal/core"
	"github.com/theandrew168/dripfile/internal/random"
)

func mockAccount(project core.Project) core.Account {
	account := core.NewAccount(
		random.String(8),
		random.String(8),
		random.String(8),
		project,
	)
	return account
}
