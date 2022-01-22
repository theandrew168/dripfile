package postgresql

import (
	"github.com/jackc/pgx/v4/pgxpool"

	"github.com/theandrew168/dripfile/internal/core"
)

type historyStorage struct {
	conn *pgxpool.Pool
}

func NewHistoryStorage(conn *pgxpool.Pool) core.HistoryStorage {
	s := historyStorage{
		conn: conn,
	}
	return &s
}

func (s *historyStorage) Create(history *core.History) error {
	return nil
}

func (s *historyStorage) Read(id int64) (core.History, error) {
	return core.History{}, nil
}

func (s *historyStorage) Update(history core.History) error {
	return nil
}

func (s *historyStorage) Delete(history core.History) error {
	return nil
}