package postgres

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v4/pgxpool"

	"github.com/theandrew168/dripfile/internal/core"
)

type transferStorage struct {
	conn *pgxpool.Pool
}

func NewTransferStorage(conn *pgxpool.Pool) core.TransferStorage {
	s := transferStorage{
		conn: conn,
	}
	return &s
}

func (s *transferStorage) Create(transfer *core.Transfer) error {
	return nil
}

func (s *transferStorage) Read(id string) (core.Transfer, error) {
	return core.Transfer{}, nil
}

func (s *transferStorage) Update(transfer core.Transfer) error {
	return nil
}

func (s *transferStorage) Delete(transfer core.Transfer) error {
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

		err := scan(rows, dest...)
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
