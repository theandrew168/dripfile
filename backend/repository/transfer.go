package repository

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"

	"github.com/theandrew168/dripfile/backend/database"
	"github.com/theandrew168/dripfile/backend/domain"
)

// ensure TransferRepository interface is satisfied
var _ TransferRepository = (*PostgresTransferRepository)(nil)

type TransferRepository interface {
	Create(transfer *domain.Transfer) error
	List() ([]*domain.Transfer, error)
	Read(id uuid.UUID) (*domain.Transfer, error)
	Update(transfer *domain.Transfer) error
	Delete(transfer *domain.Transfer) error
}

type Transfer struct {
	ID uuid.UUID `db:"id"`

	ItineraryID uuid.UUID             `db:"itinerary_id"`
	Status      domain.TransferStatus `db:"status"`
	Progress    int                   `db:"progress"`
	Error       string                `db:"error"`

	CreatedAt time.Time `db:"created_at"`
	UpdatedAt time.Time `db:"updated_at"`
}

type PostgresTransferRepository struct {
	conn database.Conn
}

func NewPostgresTransferRepository(conn database.Conn) *PostgresTransferRepository {
	repo := PostgresTransferRepository{
		conn: conn,
	}
	return &repo
}

func (repo *PostgresTransferRepository) marshal(transfer *domain.Transfer) (Transfer, error) {
	row := Transfer{
		ID: transfer.ID(),

		ItineraryID: transfer.ItineraryID(),
		Status:      transfer.Status(),
		Progress:    transfer.Progress(),
		Error:       transfer.Error(),

		CreatedAt: transfer.CreatedAt(),
		UpdatedAt: transfer.UpdatedAt(),
	}
	return row, nil
}

func (repo *PostgresTransferRepository) unmarshal(row Transfer) (*domain.Transfer, error) {
	transfer := domain.LoadTransfer(
		row.ID,
		row.ItineraryID,
		row.Status,
		row.Progress,
		row.Error,
		row.CreatedAt,
		row.UpdatedAt,
	)
	return transfer, nil
}

func (repo *PostgresTransferRepository) Create(transfer *domain.Transfer) error {
	stmt := `
		INSERT INTO transfer
			(id, itinerary_id, status, progress, error, created_at, updated_at)
		VALUES
			($1, $2, $3, $4, $5, $6, $7)`

	row, err := repo.marshal(transfer)
	if err != nil {
		return err
	}

	args := []any{
		row.ID,
		row.ItineraryID,
		row.Status,
		row.Progress,
		row.Error,
		row.CreatedAt,
		row.UpdatedAt,
	}

	ctx, cancel := context.WithTimeout(context.Background(), database.Timeout)
	defer cancel()

	_, err = repo.conn.Exec(ctx, stmt, args...)
	if err != nil {
		return checkCreateError(err)
	}

	return nil
}

func (repo *PostgresTransferRepository) List() ([]*domain.Transfer, error) {
	stmt := `
		SELECT
			id,
			itinerary_id,
			status,
			progress,
			error,
			created_at,
			updated_at
		FROM transfer
		ORDER BY created_at DESC`

	ctx, cancel := context.WithTimeout(context.Background(), database.Timeout)
	defer cancel()

	rows, err := repo.conn.Query(ctx, stmt)
	if err != nil {
		return nil, err
	}

	transferRows, err := pgx.CollectRows(rows, pgx.RowToStructByName[Transfer])
	if err != nil {
		return nil, checkListError(err)
	}

	var transfers []*domain.Transfer
	for _, row := range transferRows {
		transfer, err := repo.unmarshal(row)
		if err != nil {
			return nil, err
		}

		transfers = append(transfers, transfer)
	}

	return transfers, nil
}

func (repo *PostgresTransferRepository) Read(id uuid.UUID) (*domain.Transfer, error) {
	stmt := `
		SELECT
			id,
			itinerary_id,
			status,
			progress,
			error,
			created_at,
			updated_at
		FROM transfer
		WHERE id = $1`

	ctx, cancel := context.WithTimeout(context.Background(), database.Timeout)
	defer cancel()

	rows, err := repo.conn.Query(ctx, stmt, id)
	if err != nil {
		return nil, err
	}

	row, err := pgx.CollectOneRow(rows, pgx.RowToStructByName[Transfer])
	if err != nil {
		return nil, checkReadError(err)
	}

	return repo.unmarshal(row)
}

func (repo *PostgresTransferRepository) Update(transfer *domain.Transfer) error {
	now := time.Now()
	stmt := `
		UPDATE transfer
		SET
			status = $1,
			progress = $2,
			error = $3,
			updated_at = $4
		WHERE id = $5
		  AND updated_at = $6
		RETURNING updated_at`

	row, err := repo.marshal(transfer)
	if err != nil {
		return err
	}

	args := []any{
		row.Status,
		row.Progress,
		row.Error,
		now,
		row.ID,
		row.UpdatedAt,
	}

	ctx, cancel := context.WithTimeout(context.Background(), database.Timeout)
	defer cancel()

	rows, err := repo.conn.Query(ctx, stmt, args...)
	if err != nil {
		return err
	}

	_, err = pgx.CollectOneRow(rows, pgx.RowTo[time.Time])
	if err != nil {
		return checkUpdateError(err)
	}

	transfer.SetUpdatedAt(now)
	return err
}

func (repo *PostgresTransferRepository) Delete(transfer *domain.Transfer) error {
	stmt := `
		DELETE FROM transfer
		WHERE id = $1
		RETURNING id`

	err := transfer.CheckDelete()
	if err != nil {
		return err
	}

	ctx, cancel := context.WithTimeout(context.Background(), database.Timeout)
	defer cancel()

	rows, err := repo.conn.Query(ctx, stmt, transfer.ID())
	if err != nil {
		return err
	}

	_, err = pgx.CollectOneRow(rows, pgx.RowTo[uuid.UUID])
	if err != nil {
		return checkDeleteError(err)
	}

	return nil
}
