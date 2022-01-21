
package postgresql

import (
	"github.com/jackc/pgx/v4/pgxpool"

	"github.com/theandrew168/dripfile/internal/core"
)

type accountStorage struct {
	conn *pgxpool.Pool
}

func NewAccountStorage(conn *pgxpool.Pool) core.AccountStorage {
	s := accountStorage{
		conn: conn,
	}
	return &s
}

func (s *accountStorage) Create(account *core.Account) error {
	return nil
}

func (s *accountStorage) Read(id int64) (core.Account, error) {
	return core.Account{}, nil
}

func (s *accountStorage) Update(account core.Account) error {
	return nil
}

func (s *accountStorage) Delete(account core.Account) error {
	return nil
}
