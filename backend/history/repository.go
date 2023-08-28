package history

import (
	"context"
	"time"

	"github.com/theandrew168/dripfile/backend/database"
)

// ensure Repository interface is satisfied
var _ Repository = (*PostgresRepository)(nil)

// repository interface (other code depends on this)
type Repository interface {
	Create(l *History) error
	List() ([]*History, error)
	Read(id string) (*History, error)
}

// repository implementation (knows about domain internals)
type PostgresRepository struct {
	conn database.Conn
}

func NewRepository(conn database.Conn) *PostgresRepository {
	repo := PostgresRepository{
		conn: conn,
	}
	return &repo
}

type historyRow struct {
	id string

	totalBytes int64
	startedAt  time.Time
	finishedAt time.Time
	transferID string
}

func (repo *PostgresRepository) marshal(h *History) (historyRow, error) {
	hr := historyRow{
		id: h.ID(),

		totalBytes: h.TotalBytes(),
		startedAt:  h.StartedAt(),
		finishedAt: h.FinishedAt(),
		transferID: h.TransferID(),
	}
	return hr, nil
}

func (repo *PostgresRepository) unmarshal(hr historyRow) (*History, error) {
	return UnmarshalFromStorage(hr.id, hr.totalBytes, hr.startedAt, hr.finishedAt, hr.transferID)
}

func (repo *PostgresRepository) Create(h *History) error {
	stmt := `
		INSERT INTO history
			(id, total_bytes, started_at, finished_at, transfer_id)
		VALUES
			($1, $2, $3, $4, $5)`

	hr, err := repo.marshal(h)
	if err != nil {
		return err
	}

	args := []any{
		hr.id,
		hr.totalBytes,
		hr.startedAt,
		hr.finishedAt,
		hr.transferID,
	}

	ctx, cancel := context.WithTimeout(context.Background(), database.Timeout)
	defer cancel()

	return database.Exec(repo.conn, ctx, stmt, args...)
}

func (repo *PostgresRepository) List() ([]*History, error) {
	stmt := `
		SELECT
			id,
			total_bytes,
			started_at,
			finished_at
			transfer_id
		FROM history
		ORDER BY created_at ASC`

	ctx, cancel := context.WithTimeout(context.Background(), database.Timeout)
	defer cancel()

	rows, err := repo.conn.Query(ctx, stmt)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var hs []*History
	for rows.Next() {
		var hr historyRow
		dest := []any{
			&hr.id,
			&hr.totalBytes,
			&hr.startedAt,
			&hr.finishedAt,
			&hr.transferID,
		}

		err := database.Scan(rows, dest...)
		if err != nil {
			return nil, err
		}

		t, err := repo.unmarshal(hr)
		if err != nil {
			return nil, err
		}

		hs = append(hs, t)
	}

	if rows.Err() != nil {
		return nil, rows.Err()
	}

	return hs, nil
}

func (repo *PostgresRepository) Read(id string) (*History, error) {
	stmt := `
		SELECT
			id,
			total_bytes,
			started_at,
			finished_at
			transfer_id
		FROM history
		WHERE id = $1`

	var hr historyRow
	dest := []any{
		&hr.id,
		&hr.totalBytes,
		&hr.startedAt,
		&hr.finishedAt,
		&hr.transferID,
	}

	ctx, cancel := context.WithTimeout(context.Background(), database.Timeout)
	defer cancel()

	row := repo.conn.QueryRow(ctx, stmt, id)
	err := database.Scan(row, dest...)
	if err != nil {
		return nil, err
	}

	return repo.unmarshal(hr)
}
