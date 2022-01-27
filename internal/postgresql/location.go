package postgresql

import (
	"github.com/jackc/pgx/v4/pgxpool"

	"github.com/theandrew168/dripfile/internal/core"
)

type locationStorage struct {
	conn *pgxpool.Pool
}

func NewLocationStorage(conn *pgxpool.Pool) core.LocationStorage {
	s := locationStorage{
		conn: conn,
	}
	return &s
}

func (s *locationStorage) Create(location *core.Location) error {
	return nil
}

func (s *locationStorage) Read(id int64) (core.Location, error) {
	return core.Location{}, nil
}

func (s *locationStorage) Update(location core.Location) error {
	return nil
}

func (s *locationStorage) Delete(location core.Location) error {
	return nil
}

func (s *locationStorage) ReadAll() ([]core.Location, error) {
	return nil, nil
}
