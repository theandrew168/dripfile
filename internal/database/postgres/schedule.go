package postgres

import (
	"github.com/jackc/pgx/v4/pgxpool"

	"github.com/theandrew168/dripfile/internal/core"
)

type scheduleStorage struct {
	pool *pgxpool.Pool
}

func NewScheduleStorage(pool *pgxpool.Pool) *scheduleStorage {
	s := scheduleStorage{
		pool: pool,
	}
	return &s
}

func (s *scheduleStorage) Create(schedule *core.Schedule) error {
	return nil
}

func (s *scheduleStorage) Read(id string) (core.Schedule, error) {
	return core.Schedule{}, nil
}

func (s *scheduleStorage) Update(schedule core.Schedule) error {
	return nil
}

func (s *scheduleStorage) Delete(schedule core.Schedule) error {
	return nil
}
