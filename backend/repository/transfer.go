package repository

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"

	"github.com/theandrew168/dripfile/backend/database"
	"github.com/theandrew168/dripfile/backend/domain"
	"github.com/theandrew168/dripfile/backend/memorydb"
)

// ensure TransferRepository interface is satisfied
var _ TransferRepository = (*MemoryTransferRepository)(nil)
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
	Version   int       `db:"version"`
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
		Version:   transfer.Version(),
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
		row.Version,
	)
	return transfer, nil
}

func (repo *PostgresTransferRepository) Create(transfer *domain.Transfer) error {
	stmt := `
		INSERT INTO transfer
			(id, itinerary_id, status, progress, error)
		VALUES
			($1, $2, $3, $4, $5)`

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
			version
		FROM transfer
		ORDER BY created_at ASC`

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
			version
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
	stmt := `
		UPDATE transfer
		SET
			status = $2,
			progress = $3,
			error = $4,
			version = version + 1
		WHERE id = $1
		  AND version = $5
		RETURNING version`

	row, err := repo.marshal(transfer)
	if err != nil {
		return err
	}

	args := []any{
		row.ID,
		row.Status,
		row.Progress,
		row.Error,
		row.Version,
	}

	ctx, cancel := context.WithTimeout(context.Background(), database.Timeout)
	defer cancel()

	rows, err := repo.conn.Query(ctx, stmt, args...)
	if err != nil {
		return err
	}

	_, err = pgx.CollectOneRow(rows, pgx.RowTo[int])
	if err != nil {
		return checkUpdateError(err)
	}

	// TODO: Update domain object's Version field
	return err
}

func (repo *PostgresTransferRepository) Delete(transfer *domain.Transfer) error {
	stmt := `
		DELETE FROM transfer
		WHERE id = $1
		RETURNING version`

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

	_, err = pgx.CollectOneRow(rows, pgx.RowTo[int])
	if err != nil {
		return checkDeleteError(err)
	}

	return nil
}

type MemoryTransferRepository struct {
	db *memorydb.MemoryDB[*domain.Transfer]
}

func NewMemoryTransferRepository() *MemoryTransferRepository {
	repo := MemoryTransferRepository{
		db: memorydb.New[*domain.Transfer](),
	}
	return &repo
}

func (repo *MemoryTransferRepository) Create(transfer *domain.Transfer) error {
	return repo.db.Create(transfer)
}

func (repo *MemoryTransferRepository) List() ([]*domain.Transfer, error) {
	return repo.db.List()
}

func (repo *MemoryTransferRepository) Read(id uuid.UUID) (*domain.Transfer, error) {
	transfer, err := repo.db.Read(id)
	if err != nil {
		switch {
		case errors.Is(err, memorydb.ErrNotFound):
			return nil, ErrNotExist
		default:
			return nil, err
		}
	}

	return transfer, nil
}

func (repo *MemoryTransferRepository) Update(transfer *domain.Transfer) error {
	return repo.db.Update(transfer)
}

func (repo *MemoryTransferRepository) Delete(transfer *domain.Transfer) error {
	err := transfer.CheckDelete()
	if err != nil {
		return err
	}

	err = repo.db.Delete(transfer.ID())
	if err != nil {
		switch {
		case errors.Is(err, memorydb.ErrNotFound):
			return ErrNotExist
		default:
			return err
		}
	}

	return nil
}
