package storage

import (
	"context"
	"time"

	"github.com/theandrew168/dripfile/backend/database"
	"github.com/theandrew168/dripfile/backend/history"
)

type Storage struct {
	conn database.Conn
}

func New(conn database.Conn) *Storage {
	store := Storage{
		conn: conn,
	}
	return &store
}

type historyRow struct {
	id string

	totalBytes int64
	startedAt  time.Time
	finishedAt time.Time
	transferID string
}

func (store *Storage) marshal(h *history.History) (historyRow, error) {
	hr := historyRow{
		id: h.ID(),

		totalBytes: h.TotalBytes(),
		startedAt:  h.StartedAt(),
		finishedAt: h.FinishedAt(),
		transferID: h.TransferID(),
	}
	return hr, nil
}

func (store *Storage) unmarshal(hr historyRow) (*history.History, error) {
	return history.UnmarshalFromStorage(hr.id, hr.totalBytes, hr.startedAt, hr.finishedAt, hr.transferID)
}

func (store *Storage) Create(h *history.History) error {
	stmt := `
		INSERT INTO history
			(id, total_bytes, started_at, finished_at, transfer_id)
		VALUES
			($1, $2, $3, $4, $5)`

	hr, err := store.marshal(h)
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

	return database.Exec(store.conn, ctx, stmt, args...)
}

func (store *Storage) Read(id string) (*history.History, error) {
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

	row := store.conn.QueryRow(ctx, stmt, id)
	err := database.Scan(row, dest...)
	if err != nil {
		return nil, err
	}

	return store.unmarshal(hr)
}

func (store *Storage) List() ([]*history.History, error) {
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

	rows, err := store.conn.Query(ctx, stmt)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var hs []*history.History
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

		t, err := store.unmarshal(hr)
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
