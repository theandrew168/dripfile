package location

import (
	"context"
	"errors"

	"github.com/theandrew168/dripfile/internal/database"
)

type Repository interface {
	Create(location *Location) error
	Read(id string) (Location, error)
	List() ([]Location, error)
	Update(location Location) error
	Delete(id string) error
}

type PostgresRepository struct {
	conn database.Conn
}

func NewPostgresRepository(conn database.Conn) *PostgresRepository {
	r := PostgresRepository{
		conn: conn,
	}
	return &r
}

func (r *PostgresRepository) Create(location *Location) error {
	stmt := `
		INSERT INTO location
			(kind, name, info)
		VALUES
			($1, $2, $3)
		RETURNING id`

	args := []any{
		location.Kind,
		location.Name,
		location.Info,
	}

	ctx, cancel := context.WithTimeout(context.Background(), database.Timeout)
	defer cancel()

	row := r.conn.QueryRow(ctx, stmt, args...)
	err := database.Scan(row, &location.ID)
	if err != nil {
		if errors.Is(err, database.ErrRetry) {
			return r.Create(location)
		}

		return err
	}

	return nil
}

func (r *PostgresRepository) Read(id string) (Location, error) {
	stmt := `
		SELECT
			location.id,
			location.kind,
			location.name,
			location.info
		FROM location
		WHERE location.id = $1`

	var location Location
	dest := []any{
		&location.ID,
		&location.Kind,
		&location.Name,
		&location.Info,
	}

	ctx, cancel := context.WithTimeout(context.Background(), database.Timeout)
	defer cancel()

	row := r.conn.QueryRow(ctx, stmt, id)
	err := database.Scan(row, dest...)
	if err != nil {
		if errors.Is(err, database.ErrRetry) {
			return r.Read(id)
		}

		return Location{}, err
	}

	return location, nil
}

func (r *PostgresRepository) List() ([]Location, error) {
	stmt := `
		SELECT
			location.id,
			location.kind,
			location.name,
			location.info
		FROM location`

	ctx, cancel := context.WithTimeout(context.Background(), database.Timeout)
	defer cancel()

	rows, err := r.conn.Query(ctx, stmt)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var locations []Location
	for rows.Next() {
		var location Location
		dest := []any{
			&location.ID,
			&location.Kind,
			&location.Name,
			&location.Info,
		}

		err := database.Scan(rows, dest...)
		if err != nil {
			if errors.Is(err, database.ErrRetry) {
				return r.List()
			}

			return nil, err
		}

		locations = append(locations, location)
	}

	if rows.Err() != nil {
		return nil, rows.Err()
	}

	return locations, nil
}

func (r *PostgresRepository) Update(location Location) error {
	stmt := `
		UPDATE location
		SET
			kind = $2,
			name = $3,
			info = $4
		WHERE id = $1`

	args := []any{
		location.ID,
		location.Kind,
		location.Name,
		location.Info,
	}

	ctx, cancel := context.WithTimeout(context.Background(), database.Timeout)
	defer cancel()

	err := database.Exec(r.conn, ctx, stmt, args...)
	if err != nil {
		if errors.Is(err, database.ErrRetry) {
			return r.Update(location)
		}

		return err
	}

	return nil
}

func (r *PostgresRepository) Delete(location Location) error {
	stmt := `
		DELETE FROM location
		WHERE id = $1`

	ctx, cancel := context.WithTimeout(context.Background(), database.Timeout)
	defer cancel()

	err := database.Exec(r.conn, ctx, stmt, location.ID)
	if err != nil {
		if errors.Is(err, database.ErrRetry) {
			return r.Delete(location)
		}

		return err
	}

	return nil
}
