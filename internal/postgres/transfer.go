package postgres

import (
	"github.com/jackc/pgx/v4/pgxpool"

	"github.com/theandrew168/dripfile/internal/core"
)

type transferStorage struct {
	conn *pgxpool.Pool
}

func NewTransferStorage(conn *pgxpool.Pool) core.TransferStorage {
	s := transferStorage{
		conn: conn,
	}
	return &s
}

func (s *transferStorage) Create(transfer *core.Transfer) error {
	return nil
}

func (s *transferStorage) Read(id string) (core.Transfer, error) {
	return core.Transfer{}, nil
}

func (s *transferStorage) Update(transfer core.Transfer) error {
	return nil
}

func (s *transferStorage) Delete(transfer core.Transfer) error {
	return nil
}
