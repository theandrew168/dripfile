package memorydb

import (
	"errors"
	"sync"

	"github.com/google/uuid"
)

var ErrNotFound = errors.New("memorydb: record not found")

type Record interface {
	ID() uuid.UUID
}

type MemoryDB struct {
	sync.RWMutex
	data map[uuid.UUID]Record
}

func New() *MemoryDB {
	db := MemoryDB{
		data: make(map[uuid.UUID]Record),
	}
	return &db
}

func (db *MemoryDB) Create(record Record) error {
	db.Lock()
	defer db.Unlock()

	db.data[record.ID()] = record
	return nil
}

func (db *MemoryDB) List() ([]Record, error) {
	db.RLock()
	defer db.RUnlock()

	var records []Record
	for _, record := range db.data {
		records = append(records, record)
	}

	return records, nil
}

func (db *MemoryDB) Read(id uuid.UUID) (Record, error) {
	db.RLock()
	defer db.RUnlock()

	record, ok := db.data[id]
	if !ok {
		return nil, ErrNotFound
	}

	return record, nil
}

func (db *MemoryDB) Update(record Record) error {
	db.Lock()
	defer db.Unlock()

	db.data[record.ID()] = record
	return nil
}

func (db *MemoryDB) Delete(id uuid.UUID) error {
	db.Lock()
	defer db.Unlock()

	_, ok := db.data[id]
	if !ok {
		return ErrNotFound
	}

	delete(db.data, id)
	return nil
}
