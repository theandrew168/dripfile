package database

import (
	"github.com/theandrew168/dripfile/pkg/core"
	"github.com/theandrew168/dripfile/pkg/postgres"
)

type ScheduleStorage struct {
	db postgres.Database
}

func NewScheduleStorage(db postgres.Database) *ScheduleStorage {
	s := ScheduleStorage{
		db: db,
	}
	return &s
}

func (s *ScheduleStorage) Create(schedule *core.Schedule) error {
	return nil
}

func (s *ScheduleStorage) Read(id string) (core.Schedule, error) {
	return core.Schedule{}, nil
}

func (s *ScheduleStorage) Update(schedule core.Schedule) error {
	return nil
}

func (s *ScheduleStorage) Delete(schedule core.Schedule) error {
	return nil
}
