package postgresql

import (
	"github.com/jackc/pgx/v4/pgxpool"

	"github.com/theandrew168/dripfile/internal/core"
)

func NewStorage(conn *pgxpool.Pool) core.Storage {
	s := core.Storage{
		Account: NewAccountStorage(conn),
	}
	return s
}
