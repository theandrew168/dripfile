package storage_test

import (
	"testing"

	"github.com/theandrew168/dripfile/internal/model"
	"github.com/theandrew168/dripfile/internal/postgresql"
	"github.com/theandrew168/dripfile/internal/storage"
	"github.com/theandrew168/dripfile/internal/test"
)

func mockHistory(project model.Project) model.History {
	history := model.NewHistory(
		int64(test.RandomInt()),
		test.RandomString(8),
		test.RandomTime(),
		test.RandomTime(),
		test.RandomUUID(),
		project,
	)
	return history
}

func createHistory(t *testing.T, store *storage.Storage) (model.History, DeleterFunc) {
	t.Helper()

	project := mockProject()
	err := store.Project.Create(&project)
	test.AssertNilError(t, err)

	history := mockHistory(project)
	err = store.History.Create(&history)
	test.AssertNilError(t, err)

	deleter := func(t *testing.T) {
		err := store.History.Delete(history)
		test.AssertNilError(t, err)

		err = store.Project.Delete(project)
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

func TestHistoryDelete(t *testing.T) {
	t.Parallel()

	store, closer := test.Storage(t)
	defer closer()

	history, deleter := createHistory(t, store)
	deleter(t)

	// verify that ID isn't present anymore
	_, err := store.History.Read(history.ID)
	test.AssertErrorIs(t, err, postgresql.ErrNotExist)
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

func TestHistoryReadAllByProject(t *testing.T) {
	t.Parallel()

	store, closer := test.Storage(t)
	defer closer()

	project := mockProject()
	err := store.Project.Create(&project)
	test.AssertNilError(t, err)

	history1 := mockHistory(project)
	err = store.History.Create(&history1)
	test.AssertNilError(t, err)

	history2 := mockHistory(project)
	err = store.History.Create(&history2)
	test.AssertNilError(t, err)

	defer func(t *testing.T) {
		err := store.History.Delete(history2)
		test.AssertNilError(t, err)

		err = store.History.Delete(history1)
		test.AssertNilError(t, err)

		err = store.Project.Delete(project)
		test.AssertNilError(t, err)
	}(t)

	histories, err := store.History.ReadAllByProject(project)
	test.AssertNilError(t, err)

	var ids []string
	for _, r := range histories {
		ids = append(ids, r.ID)
	}

	test.AssertSliceContains(t, ids, history1.ID)
	test.AssertSliceContains(t, ids, history2.ID)
}
