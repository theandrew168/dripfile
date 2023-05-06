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
	repo := Repository{
		conn: conn,
	}
	return &repo
}

func (repo *Repository) Create(location *location.Location) error {
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

	row := repo.conn.QueryRow(ctx, stmt, args...)
	err := database.Scan(row, &location.ID)
	if err != nil {
		if errors.Is(err, database.ErrRetry) {
			return repo.Create(location)
		}

		return err
	}

	return nil
}

func (repo *Repository) Read(id string) (location.Location, error) {
	stmt := `
		SELECT
			id,
			kind,
			info
		FROM location
		WHERE id = $1`

	var l location.Location
	dest := []any{
		&l.ID,
		&l.Kind,
		&l.Info,
	}

	ctx, cancel := context.WithTimeout(context.Background(), database.Timeout)
	defer cancel()

	row := repo.conn.QueryRow(ctx, stmt, id)
	err := database.Scan(row, dest...)
	if err != nil {
		if errors.Is(err, database.ErrRetry) {
			return repo.Read(id)
		}

		return location.Location{}, err
	}

	return l, nil
}

func (repo *Repository) List() ([]location.Location, error) {
	stmt := `
		SELECT
			id,
			kind,
			info
		FROM location`

	ctx, cancel := context.WithTimeout(context.Background(), database.Timeout)
	defer cancel()

	rows, err := repo.conn.Query(ctx, stmt)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var ls []location.Location
	for rows.Next() {
		var l location.Location
		dest := []any{
			&l.ID,
			&l.Kind,
			&l.Info,
		}

		err := database.Scan(rows, dest...)
		if err != nil {
			if errors.Is(err, database.ErrRetry) {
				return repo.List()
			}

			return nil, err
		}

		ls = append(ls, l)
	}

	if rows.Err() != nil {
		return nil, rows.Err()
	}

	return ls, nil
}

func (repo *Repository) Update(location location.Location) error {
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
	row := repo.conn.QueryRow(ctx, stmt, args...)
	err := database.Scan(row, &updatedID)
	if err != nil {
		if errors.Is(err, database.ErrRetry) {
			return repo.Update(location)
		}

		return err
	}

	return nil
}

func (repo *Repository) Delete(id string) error {
	stmt := `
		DELETE FROM location
		WHERE id = $1
		RETURNING id`

	ctx, cancel := context.WithTimeout(context.Background(), database.Timeout)
	defer cancel()

	var deletedID string
	row := repo.conn.QueryRow(ctx, stmt, id)
	err := database.Scan(row, &deletedID)
	if err != nil {
		if errors.Is(err, database.ErrRetry) {
			return repo.Delete(id)
		}

		return err
	}

	return nil
}
