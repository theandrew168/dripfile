package storage_test

import (
	"testing"

	"github.com/theandrew168/dripfile/internal/database"
	"github.com/theandrew168/dripfile/internal/model"
	"github.com/theandrew168/dripfile/internal/storage"
	"github.com/theandrew168/dripfile/internal/test"
)

func mockAccount() model.Account {
	account := model.NewAccount(
		test.RandomString(8),
		test.RandomString(8),
		test.RandomString(8),
	)
	return account
}

func createAccount(t *testing.T, store *storage.Storage) (model.Account, DeleterFunc) {
	t.Helper()

	account := mockAccount()
	err := store.Account.Create(&account)
	test.AssertNilError(t, err)

	deleter := func(t *testing.T) {
		err := store.Account.Delete(account)
		test.AssertNilError(t, err)
	}

	return account, deleter
}

func TestAccountCreate(t *testing.T) {
	t.Parallel()

	store, closer := test.Storage(t)
	defer closer()

	account, deleter := createAccount(t, store)
	defer deleter(t)

	test.AssertNotEqual(t, account.ID, "")
}

func TestAccountCreateUnique(t *testing.T) {
	t.Parallel()

	store, closer := test.Storage(t)
	defer closer()

	account, deleter := createAccount(t, store)
	defer deleter(t)

	err := store.Account.Create(&account)
	test.AssertErrorIs(t, err, database.ErrExist)
}

func TestAccountRead(t *testing.T) {
	t.Parallel()

	store, closer := test.Storage(t)
	defer closer()

	account, deleter := createAccount(t, store)
	defer deleter(t)

	got, err := store.Account.Read(account.ID)
	test.AssertNilError(t, err)

	test.AssertEqual(t, got.ID, account.ID)
}

func TestAccountReadByEmail(t *testing.T) {
	t.Parallel()

	store, closer := test.Storage(t)
	defer closer()

	account, deleter := createAccount(t, store)
	defer deleter(t)

	got, err := store.Account.ReadByEmail(account.Email)
	test.AssertNilError(t, err)

	test.AssertEqual(t, got.ID, account.ID)
}

func TestAccountUpdate(t *testing.T) {
	t.Parallel()

	store, closer := test.Storage(t)
	defer closer()

	account, deleter := createAccount(t, store)
	defer deleter(t)

	account.Email = test.RandomString(8)
	account.Password = test.RandomString(8)
	account.Role = test.RandomString(8)
	account.Verified = true

	err := store.Account.Update(account)
	test.AssertNilError(t, err)

	got, err := store.Account.Read(account.ID)
	test.AssertNilError(t, err)

	test.AssertEqual(t, got, account)
}

func TestAccountDelete(t *testing.T) {
	t.Parallel()

	store, closer := test.Storage(t)
	defer closer()

	account, deleter := createAccount(t, store)
	deleter(t)

	// verify that the account isn't present anymore
	_, err := store.Account.Read(account.ID)
	test.AssertErrorIs(t, err, database.ErrNotExist)
}
