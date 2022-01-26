package postgresql

import (
	"github.com/jackc/pgx/v4/pgxpool"

	"github.com/theandrew168/dripfile/internal/core"
)

type jobStorage struct {
	conn *pgxpool.Pool
}

func NewJobStorage(conn *pgxpool.Pool) core.JobStorage {
	s := jobStorage{
		conn: conn,
	}
	return &s
}

func (s *jobStorage) Create(job *core.Job) error {
	return nil
}

func (s *jobStorage) Read(id int64) (core.Job, error) {
	return core.Job{}, nil
}

func (s *jobStorage) Update(job core.Job) error {
	return nil
}

func (s *jobStorage) Delete(job core.Job) error {
	return nil
}
