package pubsub

import (
	"github.com/jackc/pgx/v4/pgxpool"

	"github.com/theandrew168/dripfile/internal/database"
	"github.com/theandrew168/dripfile/internal/pubsub/postgres"
)

func NewPostgresQueue(conn *pgxpool.Pool, storage database.Storage) Queue {
	q := Queue{
		Transfer: postgres.NewTransferQueue(conn, storage),
	}
	return q
}
