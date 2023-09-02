package transfer

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/theandrew168/dripfile/backend/database"
)

// ensure Repository interface is satisfied
var _ Repository = (*PostgresRepository)(nil)

// repository interface (other code depends on this)
type Repository interface {
	Create(t *Transfer) error
	List() ([]*Transfer, error)
	Read(id uuid.UUID) (*Transfer, error)
	Delete(id uuid.UUID) error
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
	ID uuid.UUID `db:"id"`

	Pattern        string    `db:"pattern"`
	FromLocationID uuid.UUID `db:"from_location_id"`
	ToLocationID   uuid.UUID `db:"to_location_id"`

	CreatedAt time.Time `db:"created_at"`
	Version   int       `db:"version"`
}

func (repo *PostgresRepository) marshal(t *Transfer) (transferRow, error) {
	tr := transferRow{
		ID: t.id,

		Pattern:        t.pattern,
		FromLocationID: t.fromLocationID,
		ToLocationID:   t.toLocationID,

		CreatedAt: t.createdAt,
		Version:   t.version,
	}
	return tr, nil
}

func (repo *PostgresRepository) unmarshal(tr transferRow) (*Transfer, error) {
	t := Transfer{
		id: tr.ID,

		pattern:        tr.Pattern,
		fromLocationID: tr.FromLocationID,
		toLocationID:   tr.ToLocationID,

		createdAt: tr.CreatedAt,
		version:   tr.Version,
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
		tr.ID,
		tr.Pattern,
		tr.FromLocationID,
		tr.ToLocationID,
	}

	ctx, cancel := context.WithTimeout(context.Background(), database.Timeout)
	defer cancel()

	_, err = repo.conn.Exec(ctx, stmt, args...)
	if err != nil {
		return database.CheckCreateError(err)
	}

	return nil
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

	trs, err := pgx.CollectRows(rows, pgx.RowToStructByName[transferRow])
	if err != nil {
		return nil, err
	}

	var ts []*Transfer
	for _, tr := range trs {
		t, err := repo.unmarshal(tr)
		if err != nil {
			return nil, err
		}

		ts = append(ts, t)
	}

	return ts, nil
}

func (repo *PostgresRepository) Read(id uuid.UUID) (*Transfer, error) {
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

	ctx, cancel := context.WithTimeout(context.Background(), database.Timeout)
	defer cancel()

	rows, err := repo.conn.Query(ctx, stmt, id)
	if err != nil {
		return nil, err
	}

	tr, err := pgx.CollectOneRow(rows, pgx.RowToStructByName[transferRow])
	if err != nil {
		return nil, database.CheckReadError(err)
	}

	return repo.unmarshal(tr)
}

func (repo *PostgresRepository) Delete(id uuid.UUID) error {
	stmt := `
		DELETE FROM transfer
		WHERE id = $1
		RETURNING version`

	ctx, cancel := context.WithTimeout(context.Background(), database.Timeout)
	defer cancel()

	rows, err := repo.conn.Query(ctx, stmt, id)
	if err != nil {
		return err
	}

	_, err = pgx.CollectOneRow(rows, pgx.RowTo[int])
	if err != nil {
		return database.CheckDeleteError(err)
	}

	return nil
}
