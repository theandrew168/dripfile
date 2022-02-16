package database

import (
	"github.com/jackc/pgx/v4/pgxpool"

	"github.com/theandrew168/dripfile/internal/database/postgres"
)

func NewPostgresStorage(conn *pgxpool.Pool) Storage {
	s := Storage{
		Project:  postgres.NewProjectStorage(conn),
		Account:  postgres.NewAccountStorage(conn),
		Session:  postgres.NewSessionStorage(conn),
		Location: postgres.NewLocationStorage(conn),
		Transfer: postgres.NewTransferStorage(conn),
		Schedule: postgres.NewScheduleStorage(conn),
		Job:      postgres.NewJobStorage(conn),
		History:  postgres.NewHistoryStorage(conn),
	}
	return s
}
