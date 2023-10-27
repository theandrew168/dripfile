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

type MemoryDB[R Record] struct {
	sync.RWMutex
	data map[uuid.UUID]R
}

func New[R Record]() *MemoryDB[R] {
	db := MemoryDB[R]{
		data: make(map[uuid.UUID]R),
	}
	return &db
}

func (db *MemoryDB[R]) Create(record R) error {
	db.Lock()
	defer db.Unlock()

	db.data[record.ID()] = record
	return nil
}

func (db *MemoryDB[R]) List() ([]R, error) {
	db.RLock()
	defer db.RUnlock()

	var records []R
	for _, record := range db.data {
		records = append(records, record)
	}

	return records, nil
}

func (db *MemoryDB[R]) Read(id uuid.UUID) (R, error) {
	db.RLock()
	defer db.RUnlock()

	record, ok := db.data[id]
	if !ok {
		var empty R
		return empty, ErrNotFound
	}

	return record, nil
}

func (db *MemoryDB[R]) Update(record R) error {
	db.Lock()
	defer db.Unlock()

	_, ok := db.data[record.ID()]
	if !ok {
		return ErrNotFound
	}

	db.data[record.ID()] = record
	return nil
}

func (db *MemoryDB[R]) Delete(id uuid.UUID) error {
	db.Lock()
	defer db.Unlock()

	_, ok := db.data[id]
	if !ok {
		return ErrNotFound
	}

	delete(db.data, id)
	return nil
}
