package memorydb_test

import (
	"testing"

	"github.com/google/uuid"

	"github.com/theandrew168/dripfile/backend/memorydb"
	"github.com/theandrew168/dripfile/backend/test"
)

type record struct {
	id uuid.UUID
}

func newRecord() *record {
	r := record{
		id: uuid.New(),
	}
	return &r
}

func (r *record) ID() uuid.UUID {
	return r.id
}

func TestCreate(t *testing.T) {
	db := memorydb.New[*record]()

	r := newRecord()
	err := db.Create(r)
	test.AssertNilError(t, err)
}

func TestList(t *testing.T) {
	db := memorydb.New[*record]()

	r1 := newRecord()
	err := db.Create(r1)
	test.AssertNilError(t, err)

	r2 := newRecord()
	err = db.Create(r2)
	test.AssertNilError(t, err)

	rs, err := db.List()
	test.AssertNilError(t, err)
	test.AssertEqual(t, len(rs), 2)
}

func TestRead(t *testing.T) {
	db := memorydb.New[*record]()

	r := newRecord()
	err := db.Create(r)
	test.AssertNilError(t, err)

	got, err := db.Read(r.ID())
	test.AssertNilError(t, err)
	test.AssertEqual(t, got.ID(), r.ID())
}

func TestReadNotFound(t *testing.T) {
	db := memorydb.New[*record]()

	_, err := db.Read(uuid.New())
	test.AssertErrorIs(t, err, memorydb.ErrNotFound)
}

func TestUpdate(t *testing.T) {
	db := memorydb.New[*record]()

	r := newRecord()
	err := db.Create(r)
	test.AssertNilError(t, err)

	err = db.Update(r)
	test.AssertNilError(t, err)
}

func TestDelete(t *testing.T) {
	db := memorydb.New[*record]()

	r := newRecord()
	err := db.Create(r)
	test.AssertNilError(t, err)

	_, err = db.Read(r.ID())
	test.AssertNilError(t, err)

	err = db.Delete(r.ID())
	test.AssertNilError(t, err)

	_, err = db.Read(r.ID())
	test.AssertErrorIs(t, err, memorydb.ErrNotFound)
}
