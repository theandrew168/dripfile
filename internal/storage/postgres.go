package storage

import (
	"github.com/jackc/pgx/v4/pgxpool"

	"github.com/theandrew168/dripfile/internal/storage/postgres"
)

func NewPostgres(conn *pgxpool.Pool) Storage {
	s := Storage{
		Project:  postgres.NewProject(conn),
		Account:  postgres.NewAccount(conn),
		Session:  postgres.NewSession(conn),
		Location: postgres.NewLocation(conn),
		Transfer: postgres.NewTransfer(conn),
		Schedule: postgres.NewSchedule(conn),
		Job:      postgres.NewJob(conn),
		History:  postgres.NewHistory(conn),
	}
	return s
}
