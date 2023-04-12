package storage_test

import (
	"testing"

	"github.com/theandrew168/dripfile/internal/database"
	"github.com/theandrew168/dripfile/internal/model"
	"github.com/theandrew168/dripfile/internal/storage"
	"github.com/theandrew168/dripfile/internal/test"
)

func mockHistory() model.History {
	history := model.NewHistory(
		int64(test.RandomInt()),
		test.RandomString(8),
		test.RandomTime(),
		test.RandomTime(),
		test.RandomUUID(),
	)
	return history
}

func createHistory(t *testing.T, store *storage.Storage) (model.History, DeleterFunc) {
	t.Helper()

	history := mockHistory()
	err := store.History.Create(&history)
	test.AssertNilError(t, err)

	deleter := func(t *testing.T) {
		err := store.History.Delete(history)
		test.AssertNilError(t, err)
	}

	return history, deleter
}

func TestHistoryCreate(t *testing.T) {
	t.Parallel()

	store, closer := test.Storage(t)
	defer closer()

	history, deleter := createHistory(t, store)
	defer deleter(t)

	test.AssertNotEqual(t, history.ID, "")
}

func TestHistoryRead(t *testing.T) {
	t.Parallel()

	store, closer := test.Storage(t)
	defer closer()

	history, deleter := createHistory(t, store)
	defer deleter(t)

	got, err := store.History.Read(history.ID)
	test.AssertNilError(t, err)

	test.AssertEqual(t, got.ID, history.ID)
}

func TestHistoryReadAll(t *testing.T) {
	t.Parallel()

	store, closer := test.Storage(t)
	defer closer()

	history1, deleter1 := createHistory(t, store)
	defer deleter1(t)

	history2, deleter2 := createHistory(t, store)
	defer deleter2(t)

	histories, err := store.History.ReadAll()
	test.AssertNilError(t, err)

	var ids []string
	for _, r := range histories {
		ids = append(ids, r.ID)
	}

	test.AssertSliceContains(t, ids, history1.ID)
	test.AssertSliceContains(t, ids, history2.ID)
}

func TestHistoryDelete(t *testing.T) {
	t.Parallel()

	store, closer := test.Storage(t)
	defer closer()

	history, deleter := createHistory(t, store)
	deleter(t)

	// verify that the record isn't present anymore
	_, err := store.History.Read(history.ID)
	test.AssertErrorIs(t, err, database.ErrNotExist)
}
