package postgres

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v4/pgxpool"

	"github.com/theandrew168/dripfile/internal/core"
)

type memberStorage struct {
	conn *pgxpool.Pool
}

func NewMemberStorage(conn *pgxpool.Pool) core.MemberStorage {
	s := memberStorage{
		conn: conn,
	}
	return &s
}

func (s *memberStorage) Create(member *core.Member) error {
	stmt := `
		INSERT INTO member
			(role, account_id, project_id)
		VALUES
			($1, $2, $3)
		RETURNING id`

	args := []interface{}{
		member.Role,
		member.Account.ID,
		member.Project.ID,
	}

	ctx, cancel := context.WithTimeout(context.Background(), queryTimeout)
	defer cancel()

	row := s.conn.QueryRow(ctx, stmt, args...)
	err := scan(row, &member.ID)
	if err != nil {
		if errors.Is(err, core.ErrRetry) {
			return s.Create(member)
		}

		return err
	}

	return nil
}

func (s *memberStorage) Read(id string) (core.Member, error) {
	return core.Member{}, nil
}

func (s *memberStorage) Update(member core.Member) error {
	return nil
}

func (s *memberStorage) Delete(member core.Member) error {
	return nil
}

func (s *memberStorage) ReadManyByAccount(account core.Account) ([]core.Member, error) {
	return nil, nil
}

func (s *memberStorage) ReadManyByProject(project core.Project) ([]core.Member, error) {
	return nil, nil
}
