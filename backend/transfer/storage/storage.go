package storage

import (
	"context"

	"github.com/theandrew168/dripfile/backend/database"
	"github.com/theandrew168/dripfile/backend/transfer"
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

type transferRow struct {
	id string

	pattern        string
	fromLocationID string
	toLocationID   string
}

func (store *Storage) marshal(t *transfer.Transfer) (transferRow, error) {
	tr := transferRow{
		id: t.ID(),

		pattern:        t.Pattern(),
		fromLocationID: t.FromLocationID(),
		toLocationID:   t.ToLocationID(),
	}
	return tr, nil
}

func (store *Storage) unmarshal(tr transferRow) (*transfer.Transfer, error) {
	return transfer.UnmarshalFromStorage(tr.id, tr.pattern, tr.fromLocationID, tr.toLocationID)
}

func (store *Storage) Create(t *transfer.Transfer) error {
	stmt := `
		INSERT INTO transfer
			(id, pattern, from_location_id, to_location_id)
		VALUES
			($1, $2, $3, $4)`

	tr, err := store.marshal(t)
	if err != nil {
		return err
	}

	args := []any{
		tr.id,
		tr.pattern,
		tr.fromLocationID,
		tr.toLocationID,
	}

	ctx, cancel := context.WithTimeout(context.Background(), database.Timeout)
	defer cancel()

	return database.Exec(store.conn, ctx, stmt, args...)
}

func (store *Storage) Read(id string) (*transfer.Transfer, error) {
	stmt := `
		SELECT
			id,
			pattern,
			from_location_id,
			to_location_id
		FROM transfer
		WHERE id = $1`

	var tr transferRow
	dest := []any{
		&tr.id,
		&tr.pattern,
		&tr.fromLocationID,
		&tr.toLocationID,
	}

	ctx, cancel := context.WithTimeout(context.Background(), database.Timeout)
	defer cancel()

	row := store.conn.QueryRow(ctx, stmt, id)
	err := database.Scan(row, dest...)
	if err != nil {
		return nil, err
	}

	return store.unmarshal(tr)
}

func (store *Storage) List() ([]*transfer.Transfer, error) {
	stmt := `
		SELECT
			id,
			pattern,
			from_location_id,
			to_location_id
		FROM transfer
		ORDER BY created_at ASC`

	ctx, cancel := context.WithTimeout(context.Background(), database.Timeout)
	defer cancel()

	rows, err := store.conn.Query(ctx, stmt)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var ts []*transfer.Transfer
	for rows.Next() {
		var tr transferRow
		dest := []any{
			&tr.id,
			&tr.pattern,
			&tr.fromLocationID,
			&tr.toLocationID,
		}

		err := database.Scan(rows, dest...)
		if err != nil {
			return nil, err
		}

		t, err := store.unmarshal(tr)
		if err != nil {
			return nil, err
		}

		ts = append(ts, t)
	}

	if rows.Err() != nil {
		return nil, rows.Err()
	}

	return ts, nil
}

func (store *Storage) Delete(id string) error {
	stmt := `
		DELETE FROM transfer
		WHERE id = $1
		RETURNING id`

	ctx, cancel := context.WithTimeout(context.Background(), database.Timeout)
	defer cancel()

	var deletedID string
	row := store.conn.QueryRow(ctx, stmt, id)
	return database.Scan(row, &deletedID)
}
