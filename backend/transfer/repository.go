package transfer

import (
	"context"
	"time"

	"github.com/theandrew168/dripfile/backend/database"
)

// ensure Repository interface is satisfied
var _ Repository = (*PostgresRepository)(nil)

// repository interface (other code depends on this)
type Repository interface {
	Create(t *Transfer) error
	List() ([]*Transfer, error)
	Read(id string) (*Transfer, error)
	Delete(id string) error
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

type transferRow struct {
	id string

	pattern        string
	fromLocationID string
	toLocationID   string

	createdAt time.Time
	version   int
}

func (repo *PostgresRepository) marshal(t *Transfer) (transferRow, error) {
	tr := transferRow{
		id: t.id,

		pattern:        t.pattern,
		fromLocationID: t.fromLocationID,
		toLocationID:   t.toLocationID,

		createdAt: t.createdAt,
		version:   t.version,
	}
	return tr, nil
}

func (repo *PostgresRepository) unmarshal(tr transferRow) (*Transfer, error) {
	t := Transfer{
		id: tr.id,

		pattern:        tr.pattern,
		fromLocationID: tr.fromLocationID,
		toLocationID:   tr.toLocationID,

		createdAt: tr.createdAt,
		version:   tr.version,
	}
	return &t, nil
}

func (repo *PostgresRepository) Create(t *Transfer) error {
	stmt := `
		INSERT INTO transfer
			(id, pattern, from_location_id, to_location_id)
		VALUES
			($1, $2, $3, $4)`

	tr, err := repo.marshal(t)
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

	return database.Exec(repo.conn, ctx, stmt, args...)
}

func (repo *PostgresRepository) List() ([]*Transfer, error) {
	stmt := `
		SELECT
			id,
			pattern,
			from_location_id,
			to_location_id,
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
	defer rows.Close()

	var ts []*Transfer
	for rows.Next() {
		var tr transferRow
		dest := []any{
			&tr.id,
			&tr.pattern,
			&tr.fromLocationID,
			&tr.toLocationID,
			&tr.createdAt,
			&tr.version,
		}

		err := database.Scan(rows, dest...)
		if err != nil {
			return nil, err
		}

		t, err := repo.unmarshal(tr)
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

func (repo *PostgresRepository) Read(id string) (*Transfer, error) {
	stmt := `
		SELECT
			id,
			pattern,
			from_location_id,
			to_location_id,
			created_at,
			version
		FROM transfer
		WHERE id = $1`

	var tr transferRow
	dest := []any{
		&tr.id,
		&tr.pattern,
		&tr.fromLocationID,
		&tr.toLocationID,
		&tr.createdAt,
		&tr.version,
	}

	ctx, cancel := context.WithTimeout(context.Background(), database.Timeout)
	defer cancel()

	row := repo.conn.QueryRow(ctx, stmt, id)
	err := database.Scan(row, dest...)
	if err != nil {
		return nil, err
	}

	return repo.unmarshal(tr)
}

func (repo *PostgresRepository) Delete(id string) error {
	stmt := `
		DELETE FROM transfer
		WHERE id = $1
		RETURNING id`

	ctx, cancel := context.WithTimeout(context.Background(), database.Timeout)
	defer cancel()

	var deletedID string
	row := repo.conn.QueryRow(ctx, stmt, id)
	return database.Scan(row, &deletedID)
}
