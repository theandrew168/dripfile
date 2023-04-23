package repository

import (
	"context"
	"errors"

	"github.com/theandrew168/dripfile/internal/database"
	"github.com/theandrew168/dripfile/internal/location"
)

type Repository struct {
	conn database.Conn
}

func New(conn database.Conn) *Repository {
	r := Repository{
		conn: conn,
	}
	return &r
}

func (r *Repository) Create(location *location.Location) error {
	stmt := `
		INSERT INTO location
			(kind, info)
		VALUES
			($1, $2)
		RETURNING id`

	args := []any{
		location.Kind,
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

func (r *Repository) Read(id string) (location.Location, error) {
	stmt := `
		SELECT
			location.id,
			location.kind,
			location.info
		FROM location
		WHERE location.id = $1`

	var m location.Location
	dest := []any{
		&m.ID,
		&m.Kind,
		&m.Info,
	}

	ctx, cancel := context.WithTimeout(context.Background(), database.Timeout)
	defer cancel()

	row := r.conn.QueryRow(ctx, stmt, id)
	err := database.Scan(row, dest...)
	if err != nil {
		if errors.Is(err, database.ErrRetry) {
			return r.Read(id)
		}

		return location.Location{}, err
	}

	return m, nil
}

func (r *Repository) List() ([]location.Location, error) {
	stmt := `
		SELECT
			location.id,
			location.kind,
			location.info
		FROM location`

	ctx, cancel := context.WithTimeout(context.Background(), database.Timeout)
	defer cancel()

	rows, err := r.conn.Query(ctx, stmt)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var ms []location.Location
	for rows.Next() {
		var m location.Location
		dest := []any{
			&m.ID,
			&m.Kind,
			&m.Info,
		}

		err := database.Scan(rows, dest...)
		if err != nil {
			if errors.Is(err, database.ErrRetry) {
				return r.List()
			}

			return nil, err
		}

		ms = append(ms, m)
	}

	if rows.Err() != nil {
		return nil, rows.Err()
	}

	return ms, nil
}

func (r *Repository) Update(location location.Location) error {
	stmt := `
		UPDATE location
		SET
			kind = $2,
			info = $3
		WHERE id = $1
		RETURNING id`

	args := []any{
		location.ID,
		location.Kind,
		location.Info,
	}

	ctx, cancel := context.WithTimeout(context.Background(), database.Timeout)
	defer cancel()

	var updatedID string
	row := r.conn.QueryRow(ctx, stmt, args...)
	err := database.Scan(row, &updatedID)
	if err != nil {
		if errors.Is(err, database.ErrRetry) {
			return r.Update(location)
		}

		return err
	}

	return nil
}

func (r *Repository) Delete(id string) error {
	stmt := `
		DELETE FROM location
		WHERE id = $1
		RETURNING id`

	ctx, cancel := context.WithTimeout(context.Background(), database.Timeout)
	defer cancel()

	var deletedID string
	row := r.conn.QueryRow(ctx, stmt, id)
	err := database.Scan(row, &deletedID)
	if err != nil {
		if errors.Is(err, database.ErrRetry) {
			return r.Delete(id)
		}

		return err
	}

	return nil
}
