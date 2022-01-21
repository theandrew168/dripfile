package postgresql

import (
	"github.com/jackc/pgx/v4/pgxpool"

	"github.com/theandrew168/dripfile/internal/core"
)

type scheduleStorage struct {
	conn *pgxpool.Pool
}

func NewScheduleStorage(conn *pgxpool.Pool) core.ScheduleStorage {
	s := scheduleStorage{
		conn: conn,
	}
	return &s
}

func (s *scheduleStorage) Create(schedule *core.Schedule) error {
	return nil
}

func (s *scheduleStorage) Read(id int64) (core.Schedule, error) {
	return core.Schedule{}, nil
}

func (s *scheduleStorage) Update(schedule core.Schedule) error {
	return nil
}

func (s *scheduleStorage) Delete(schedule core.Schedule) error {
	return nil
}
