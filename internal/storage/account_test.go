package storage_test

import (
	"testing"

	"github.com/theandrew168/dripfile/internal/model"
	"github.com/theandrew168/dripfile/internal/postgresql"
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

func TestAccount(t *testing.T) {
	store, closer := test.Storage(t)
	defer closer()

	// before
	project := mockProject()
	err := store.Project.Create(&project)
	if err != nil {
		t.Fatal(err)
	}

	// create
	account := mockAccount(project)
	err = store.Account.Create(&account)
	test.AssertNilError(t, err)

	test.AssertNotEqual(t, account.ID, "")

	// create (duplicate)
	err = store.Account.Create(&account)
	test.AssertErrorIs(t, err, postgresql.ErrExist)

	// read
	got, err := store.Account.Read(account.ID)
	test.AssertNilError(t, err)

	test.AssertEqual(t, got.ID, account.ID)

	// delete
	err = store.Account.Delete(account)
	test.AssertNilError(t, err)

	// verify that ID isn't present anymore
	_, err = store.Account.Read(account.ID)
	test.AssertErrorIs(t, err, postgresql.ErrNotExist)

	// after
	err = store.Project.Delete(project)
	test.AssertNilError(t, err)
}
