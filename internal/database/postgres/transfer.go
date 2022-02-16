package postgres

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v4/pgxpool"

	"github.com/theandrew168/dripfile/internal/core"
	"github.com/theandrew168/dripfile/internal/postgres"
)

type transferStorage struct {
	conn *pgxpool.Pool
}

func NewTransferStorage(conn *pgxpool.Pool) *transferStorage {
	s := transferStorage{
		conn: conn,
	}
	return &s
}

func (s *transferStorage) Create(transfer *core.Transfer) error {
	stmt := `
		INSERT INTO transfer
			(pattern, src_id, dst_id, project_id)
		VALUES
			($1, $2, $3, $4)
		RETURNING id`

	args := []interface{}{
		transfer.Pattern,
		transfer.Src.ID,
		transfer.Dst.ID,
		transfer.Project.ID,
	}

	ctx, cancel := context.WithTimeout(context.Background(), queryTimeout)
	defer cancel()

	row := s.conn.QueryRow(ctx, stmt, args...)
	err := postgres.Scan(row, &transfer.ID)
	if err != nil {
		if errors.Is(err, core.ErrRetry) {
			return s.Create(transfer)
		}

		return err
	}

	return nil
}

func (s *transferStorage) Read(id string) (core.Transfer, error) {
	stmt := `
		SELECT
			transfer.id,
			transfer.pattern,
			src.id,
			src.kind,
			src.info,
			src.project_id,
			dst.id,
			dst.kind,
			dst.info,
			dst.project_id,
			project.id
		FROM transfer
		INNER JOIN location src
			ON src.id = transfer.src_id
		INNER JOIN location dst
			ON dst.id = transfer.dst_id
		INNER JOIN project
			ON project.id = transfer.project_id
		WHERE transfer.id = $1`

	var transfer core.Transfer
	dest := []interface{}{
		&transfer.ID,
		&transfer.Pattern,
		&transfer.Src.ID,
		&transfer.Src.Kind,
		&transfer.Src.Info,
		&transfer.Src.Project.ID,
		&transfer.Dst.ID,
		&transfer.Dst.Kind,
		&transfer.Dst.Info,
		&transfer.Dst.Project.ID,
		&transfer.Project.ID,
	}

	ctx, cancel := context.WithTimeout(context.Background(), queryTimeout)
	defer cancel()

	row := s.conn.QueryRow(ctx, stmt, id)
	err := postgres.Scan(row, dest...)
	if err != nil {
		if errors.Is(err, core.ErrRetry) {
			return s.Read(id)
		}

		return core.Transfer{}, err
	}

	return transfer, nil
}

func (s *transferStorage) Update(transfer core.Transfer) error {
	stmt := `
		UPDATE transfer
		SET
			pattern = $2,
			src_id = $3,
			dst_id = $4
		WHERE id = $1`

	args := []interface{}{
		transfer.ID,
		transfer.Pattern,
		transfer.Src.ID,
		transfer.Dst.ID,
	}

	ctx, cancel := context.WithTimeout(context.Background(), queryTimeout)
	defer cancel()

	err := postgres.Exec(s.conn, ctx, stmt, args...)
	if err != nil {
		if errors.Is(err, core.ErrRetry) {
			return s.Update(transfer)
		}

		return err
	}

	return nil
}

func (s *transferStorage) Delete(transfer core.Transfer) error {
	stmt := `
		DELETE FROM transfer
		WHERE id = $1`

	ctx, cancel := context.WithTimeout(context.Background(), queryTimeout)
	defer cancel()

	err := postgres.Exec(s.conn, ctx, stmt, transfer.ID)
	if err != nil {
		if errors.Is(err, core.ErrRetry) {
			return s.Delete(transfer)
		}

		return err
	}

	return nil
}

func (s *transferStorage) ReadManyByProject(project core.Project) ([]core.Transfer, error) {
	stmt := `
		SELECT
			transfer.id,
			transfer.pattern,
			src.id,
			src.kind,
			src.info,
			src.project_id,
			dst.id,
			dst.kind,
			dst.info,
			dst.project_id,
			project.id
		FROM transfer
		INNER JOIN location src
			ON src.id = transfer.src_id
		INNER JOIN location dst
			ON dst.id = transfer.dst_id
		INNER JOIN project
			ON project.id = transfer.project_id
		WHERE project.id = $1`

	ctx, cancel := context.WithTimeout(context.Background(), queryTimeout)
	defer cancel()

	rows, err := s.conn.Query(ctx, stmt, project.ID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var transfers []core.Transfer
	for rows.Next() {
		var transfer core.Transfer
		dest := []interface{}{
			&transfer.ID,
			&transfer.Pattern,
			&transfer.Src.ID,
			&transfer.Src.Kind,
			&transfer.Src.Info,
			&transfer.Src.Project.ID,
			&transfer.Dst.ID,
			&transfer.Dst.Kind,
			&transfer.Dst.Info,
			&transfer.Dst.Project.ID,
			&transfer.Project.ID,
		}

		err := postgres.Scan(rows, dest...)
		if err != nil {
			if errors.Is(err, core.ErrRetry) {
				return s.ReadManyByProject(project)
			}

			return nil, err
		}

		transfers = append(transfers, transfer)
	}

	if rows.Err() != nil {
		return nil, rows.Err()
	}

	return transfers, nil
}
