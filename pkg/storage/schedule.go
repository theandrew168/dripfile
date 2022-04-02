package storage

import (
	"github.com/theandrew168/dripfile/pkg/core"
	"github.com/theandrew168/dripfile/pkg/postgres"
)

type Schedule struct {
	pg postgres.Interface
}

func NewSchedule(pg postgres.Interface) *Schedule {
	s := Schedule{
		pg: pg,
	}
	return &s
}

func (s *Schedule) Create(schedule *core.Schedule) error {
	return nil
}

func (s *Schedule) Read(id string) (core.Schedule, error) {
	return core.Schedule{}, nil
}

func (s *Schedule) Update(schedule core.Schedule) error {
	return nil
}

func (s *Schedule) Delete(schedule core.Schedule) error {
	return nil
}
