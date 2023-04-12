package storage_test

import (
	"github.com/theandrew168/dripfile/internal/database"
	"github.com/theandrew168/dripfile/internal/model"
	"github.com/theandrew168/dripfile/internal/storage"
	"github.com/theandrew168/dripfile/internal/test"

	"testing"
)

func mockLocation() model.Location {
	location := model.NewLocation(
		test.RandomString(8),
		test.RandomString(8),
		test.RandomBytes(8),
	)
	return location
}

func createLocation(t *testing.T, store *storage.Storage) (model.Location, DeleterFunc) {
	t.Helper()

	location := mockLocation()
	err := store.Location.Create(&location)
	test.AssertNilError(t, err)

	deleter := func(t *testing.T) {
		err := store.Location.Delete(location)
		test.AssertNilError(t, err)
	}

	return location, deleter
}

func TestLocationCreate(t *testing.T) {
	t.Parallel()

	store, closer := test.Storage(t)
	defer closer()

	location, deleter := createLocation(t, store)
	defer deleter(t)

	test.AssertNotEqual(t, location.ID, "")
}

func TestLocationRead(t *testing.T) {
	t.Parallel()

	store, closer := test.Storage(t)
	defer closer()

	location, deleter := createLocation(t, store)
	defer deleter(t)

	got, err := store.Location.Read(location.ID)
	test.AssertNilError(t, err)

	test.AssertEqual(t, got.ID, location.ID)
}

func TestLocationReadAll(t *testing.T) {
	t.Parallel()

	store, closer := test.Storage(t)
	defer closer()

	location1, deleter1 := createLocation(t, store)
	defer deleter1(t)

	location2, deleter2 := createLocation(t, store)
	defer deleter2(t)

	locations, err := store.Location.ReadAll()
	test.AssertNilError(t, err)

	var ids []string
	for _, r := range locations {
		ids = append(ids, r.ID)
	}

	test.AssertSliceContains(t, ids, location1.ID)
	test.AssertSliceContains(t, ids, location2.ID)
}

func TestLocationUpdate(t *testing.T) {
	t.Parallel()

	store, closer := test.Storage(t)
	defer closer()

	location, deleter := createLocation(t, store)
	defer deleter(t)

	location.Kind = test.RandomString(8)
	location.Name = test.RandomString(8)
	location.Info = test.RandomBytes(8)

	err := store.Location.Update(location)
	test.AssertNilError(t, err)

	got, err := store.Location.Read(location.ID)
	test.AssertNilError(t, err)

	test.AssertEqual(t, got, location)
}

func TestLocationDelete(t *testing.T) {
	t.Parallel()

	store, closer := test.Storage(t)
	defer closer()

	location, deleter := createLocation(t, store)
	deleter(t)

	// verify that the record isn't present anymore
	_, err := store.Location.Read(location.ID)
	test.AssertErrorIs(t, err, database.ErrNotExist)
}
