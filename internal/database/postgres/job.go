package postgres

import (
	"github.com/jackc/pgx/v4/pgxpool"

	"github.com/theandrew168/dripfile/internal/core"
)

type jobStorage struct {
	pool *pgxpool.Pool
}

func NewJobStorage(pool *pgxpool.Pool) *jobStorage {
	s := jobStorage{
		pool: pool,
	}
	return &s
}

func (s *jobStorage) Create(job *core.Job) error {
	return nil
}

func (s *jobStorage) Read(id string) (core.Job, error) {
	return core.Job{}, nil
}

func (s *jobStorage) Update(job core.Job) error {
	return nil
}

func (s *jobStorage) Delete(job core.Job) error {
	return nil
}
