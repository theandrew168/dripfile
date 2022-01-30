package postgres

import (
	"github.com/jackc/pgx/v4/pgxpool"

	"github.com/theandrew168/dripfile/internal/core"
)

func NewStorage(conn *pgxpool.Pool) core.Storage {
	s := core.Storage{
		Project:  NewProjectStorage(conn),
		Account:  NewAccountStorage(conn),
		Session:  NewSessionStorage(conn),
		Location: NewLocationStorage(conn),
		Transfer: NewTransferStorage(conn),
		Schedule: NewScheduleStorage(conn),
		Job:      NewJobStorage(conn),
		History:  NewHistoryStorage(conn),
	}
	return s
}
