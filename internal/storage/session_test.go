package storage_test

import (
	"github.com/theandrew168/dripfile/internal/model"
	"github.com/theandrew168/dripfile/internal/postgresql"
	"github.com/theandrew168/dripfile/internal/storage"
	"github.com/theandrew168/dripfile/internal/test"

	"testing"
	"time"
)

func mockSession(account model.Account) model.Session {
	session := model.NewSession(
		test.RandomString(8),
		test.RandomTime(),
		account,
	)
	return session
}

func createSession(t *testing.T, store *storage.Storage) (model.Session, DeleterFunc) {
	t.Helper()

	account := mockAccount()
	err := store.Account.Create(&account)
	test.AssertNilError(t, err)

	session := mockSession(account)
	err = store.Session.Create(&session)
	test.AssertNilError(t, err)

	deleter := func(t *testing.T) {
		err := store.Session.Delete(session)
		test.AssertNilError(t, err)

		err = store.Account.Delete(account)
		test.AssertNilError(t, err)
	}

	return session, deleter
}

func TestSessionCreate(t *testing.T) {
	t.Parallel()

	store, closer := test.Storage(t)
	defer closer()

	session, deleter := createSession(t, store)
	defer deleter(t)

	test.AssertNotEqual(t, session.Hash, "")
}

func TestSessionRead(t *testing.T) {
	t.Parallel()

	store, closer := test.Storage(t)
	defer closer()

	session, deleter := createSession(t, store)
	defer deleter(t)

	got, err := store.Session.Read(session.Hash)
	test.AssertNilError(t, err)

	test.AssertEqual(t, got.Hash, session.Hash)
}

func TestSessionDelete(t *testing.T) {
	t.Parallel()

	store, closer := test.Storage(t)
	defer closer()

	session, deleter := createSession(t, store)
	deleter(t)

	// verify that the record isn't present anymore
	_, err := store.Session.Read(session.Hash)
	test.AssertErrorIs(t, err, postgresql.ErrNotExist)
}

func TestSessionDeleteExpired(t *testing.T) {
	t.Parallel()

	store, closer := test.Storage(t)
	defer closer()

	// mock session's expiry is now (should get deleted)
	session, deleter := createSession(t, store)
	defer deleter(t)

	offset := 5 * time.Minute

	// create an older session (should get deleted)
	older := model.NewSession(
		test.RandomString(8),
		session.Expiry.Add(-offset),
		session.Account,
	)
	err := store.Session.Create(&older)
	test.AssertNilError(t, err)

	// create a newer session (should NOT get deleted)
	newer := model.NewSession(
		test.RandomString(8),
		session.Expiry.Add(offset),
		session.Account,
	)
	err = store.Session.Create(&newer)
	test.AssertNilError(t, err)

	err = store.Session.DeleteExpired()
	test.AssertNilError(t, err)

	// verify that the older sessions got deleted
	_, err = store.Session.Read(session.Hash)
	test.AssertErrorIs(t, err, postgresql.ErrNotExist)
	_, err = store.Session.Read(older.Hash)
	test.AssertErrorIs(t, err, postgresql.ErrNotExist)

	// verify that the newer session did NOT get deleted
	_, err = store.Session.Read(newer.Hash)
	test.AssertNilError(t, err)
}
