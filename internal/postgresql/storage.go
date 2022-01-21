package postgresql

import (
	"github.com/jackc/pgx/v4/pgxpool"

	"github.com/theandrew168/dripfile/internal/core"
)

func NewStorage(conn *pgxpool.Pool) core.Storage {
	s := core.Storage{
		Account:  NewAccountStorage(conn),
		Session:  NewSessionStorage(conn),
		Project:  NewProjectStorage(conn),
		Location: NewLocationStorage(conn),
		Transfer: NewTransferStorage(conn),
		Schedule: NewScheduleStorage(conn),
		History:  NewHistoryStorage(conn),
	}
	return s
}
