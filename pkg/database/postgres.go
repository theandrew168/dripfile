package database

import (
	"github.com/jackc/pgx/v4/pgxpool"

	"github.com/theandrew168/dripfile/pkg/database/postgres"
)

func NewPostgresStorage(pool *pgxpool.Pool) Storage {
	s := Storage{
		Project:  postgres.NewProjectStorage(pool),
		Account:  postgres.NewAccountStorage(pool),
		Session:  postgres.NewSessionStorage(pool),
		Location: postgres.NewLocationStorage(pool),
		Transfer: postgres.NewTransferStorage(pool),
		Schedule: postgres.NewScheduleStorage(pool),
		Job:      postgres.NewJobStorage(pool),
		History:  postgres.NewHistoryStorage(pool),
	}
	return s
}
