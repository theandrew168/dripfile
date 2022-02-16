package postgres

import (
	"github.com/jackc/pgx/v4/pgxpool"

	"github.com/theandrew168/dripfile/internal/core"
)

type transferQueue struct {
	conn *pgxpool.Pool
}

func NewTransferQueue(conn *pgxpool.Pool) *transferQueue {
	q := transferQueue{
		conn: conn,
	}
	return &q
}

func (q *transferQueue) Push(transfer core.Transfer) error {
	return nil
}

func (q *transferQueue) Pop() (core.Transfer, error) {
	return core.Transfer{}, nil
}
