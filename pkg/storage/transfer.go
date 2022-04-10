package storage

import (
	"context"
	"errors"

	"github.com/theandrew168/dripfile/pkg/core"
	"github.com/theandrew168/dripfile/pkg/postgres"
)

type Transfer struct {
	pg postgres.Interface
}

func NewTransfer(pg postgres.Interface) *Transfer {
	s := Transfer{
		pg: pg,
	}
	return &s
}

func (s *Transfer) Create(transfer *core.Transfer) error {
	stmt := `
		INSERT INTO transfer
			(pattern, src_id, dst_id, schedule_id, project_id)
		VALUES
			($1, $2, $3, $4, $5)
		RETURNING id`

	args := []interface{}{
		transfer.Pattern,
		transfer.Src.ID,
		transfer.Dst.ID,
		transfer.Schedule.ID,
		transfer.Project.ID,
	}

	ctx, cancel := context.WithTimeout(context.Background(), queryTimeout)
	defer cancel()

	row := s.pg.QueryRow(ctx, stmt, args...)
	err := postgres.Scan(row, &transfer.ID)
	if err != nil {
		if errors.Is(err, core.ErrRetry) {
			return s.Create(transfer)
		}

		return err
	}

	return nil
}

func (s *Transfer) Read(id string) (core.Transfer, error) {
	stmt := `
		SELECT
			transfer.id,
			transfer.pattern,
			src.id,
			src.kind,
			src.name,
			src.info,
			src.project_id,
			dst.id,
			dst.kind,
			dst.name,
			dst.info,
			dst.project_id,
			schedule.id,
			schedule.name,
			schedule.expr,
			schedule.project_id,
			project.id,
			project.customer_id,
			project.subscription_item_id
		FROM transfer
		INNER JOIN location src
			ON src.id = transfer.src_id
		INNER JOIN location dst
			ON dst.id = transfer.dst_id
		INNER JOIN schedule
			ON schedule.id = transfer.schedule_id
		INNER JOIN project
			ON project.id = transfer.project_id
		WHERE transfer.id = $1`

	var transfer core.Transfer
	dest := []interface{}{
		&transfer.ID,
		&transfer.Pattern,
		&transfer.Src.ID,
		&transfer.Src.Kind,
		&transfer.Src.Name,
		&transfer.Src.Info,
		&transfer.Src.Project.ID,
		&transfer.Dst.ID,
		&transfer.Dst.Kind,
		&transfer.Dst.Name,
		&transfer.Dst.Info,
		&transfer.Dst.Project.ID,
		&transfer.Schedule.ID,
		&transfer.Schedule.Name,
		&transfer.Schedule.Expr,
		&transfer.Schedule.Project.ID,
		&transfer.Project.ID,
		&transfer.Project.CustomerID,
		&transfer.Project.SubscriptionItemID,
	}

	ctx, cancel := context.WithTimeout(context.Background(), queryTimeout)
	defer cancel()

	row := s.pg.QueryRow(ctx, stmt, id)
	err := postgres.Scan(row, dest...)
	if err != nil {
		if errors.Is(err, core.ErrRetry) {
			return s.Read(id)
		}

		return core.Transfer{}, err
	}

	return transfer, nil
}

func (s *Transfer) Update(transfer core.Transfer) error {
	stmt := `
		UPDATE transfer
		SET
			pattern = $2,
			src_id = $3,
			dst_id = $4,
			schedule_id = $5,
		WHERE id = $1`

	args := []interface{}{
		transfer.ID,
		transfer.Pattern,
		transfer.Src.ID,
		transfer.Dst.ID,
		transfer.Schedule.ID,
	}

	ctx, cancel := context.WithTimeout(context.Background(), queryTimeout)
	defer cancel()

	err := postgres.Exec(s.pg, ctx, stmt, args...)
	if err != nil {
		if errors.Is(err, core.ErrRetry) {
			return s.Update(transfer)
		}

		return err
	}

	return nil
}

func (s *Transfer) Delete(transfer core.Transfer) error {
	stmt := `
		DELETE FROM transfer
		WHERE id = $1`

	ctx, cancel := context.WithTimeout(context.Background(), queryTimeout)
	defer cancel()

	err := postgres.Exec(s.pg, ctx, stmt, transfer.ID)
	if err != nil {
		if errors.Is(err, core.ErrRetry) {
			return s.Delete(transfer)
		}

		return err
	}

	return nil
}

func (s *Transfer) ReadAllByProject(project core.Project) ([]core.Transfer, error) {
	stmt := `
		SELECT
			transfer.id,
			transfer.pattern,
			src.id,
			src.kind,
			src.name,
			src.info,
			src.project_id,
			dst.id,
			dst.kind,
			dst.name,
			dst.info,
			dst.project_id,
			schedule.id,
			schedule.name,
			schedule.expr,
			schedule.project_id,
			project.id,
			project.customer_id,
			project.subscription_item_id
		FROM transfer
		INNER JOIN location src
			ON src.id = transfer.src_id
		INNER JOIN location dst
			ON dst.id = transfer.dst_id
		INNER JOIN schedule
			ON schedule.id = transfer.schedule_id
		INNER JOIN project
			ON project.id = transfer.project_id
		WHERE project.id = $1`

	ctx, cancel := context.WithTimeout(context.Background(), queryTimeout)
	defer cancel()

	rows, err := s.pg.Query(ctx, stmt, project.ID)
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
			&transfer.Src.Name,
			&transfer.Src.Info,
			&transfer.Src.Project.ID,
			&transfer.Dst.ID,
			&transfer.Dst.Kind,
			&transfer.Dst.Name,
			&transfer.Dst.Info,
			&transfer.Dst.Project.ID,
			&transfer.Schedule.ID,
			&transfer.Schedule.Name,
			&transfer.Schedule.Expr,
			&transfer.Schedule.Project.ID,
			&transfer.Project.ID,
			&transfer.Project.CustomerID,
			&transfer.Project.SubscriptionItemID,
		}

		err := postgres.Scan(rows, dest...)
		if err != nil {
			if errors.Is(err, core.ErrRetry) {
				return s.ReadAllByProject(project)
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
