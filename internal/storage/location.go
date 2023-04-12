package storage

import (
	"context"
	"errors"

	"github.com/theandrew168/dripfile/internal/database"
	"github.com/theandrew168/dripfile/internal/model"
)

type Location struct {
	db database.Conn
}

func NewLocation(db database.Conn) *Location {
	s := Location{
		db: db,
	}
	return &s
}

func (s *Location) Create(location *model.Location) error {
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

	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	row := s.db.QueryRow(ctx, stmt, args...)
	err := database.Scan(row, &location.ID)
	if err != nil {
		if errors.Is(err, database.ErrRetry) {
			return s.Create(location)
		}

		return err
	}

	return nil
}

func (s *Location) Read(id string) (model.Location, error) {
	stmt := `
		SELECT
			location.id,
			location.kind,
			location.name,
			location.info
		FROM location
		WHERE location.id = $1`

	var location model.Location
	dest := []any{
		&location.ID,
		&location.Kind,
		&location.Name,
		&location.Info,
	}

	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	row := s.db.QueryRow(ctx, stmt, id)
	err := database.Scan(row, dest...)
	if err != nil {
		if errors.Is(err, database.ErrRetry) {
			return s.Read(id)
		}

		return model.Location{}, err
	}

	return location, nil
}

func (s *Location) Update(location model.Location) error {
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

	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	err := database.Exec(s.db, ctx, stmt, args...)
	if err != nil {
		if errors.Is(err, database.ErrRetry) {
			return s.Update(location)
		}

		return err
	}

	return nil
}

func (s *Location) Delete(location model.Location) error {
	stmt := `
		DELETE FROM location
		WHERE id = $1`

	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	err := database.Exec(s.db, ctx, stmt, location.ID)
	if err != nil {
		if errors.Is(err, database.ErrRetry) {
			return s.Delete(location)
		}

		return err
	}

	return nil
}

func (s *Location) ReadAll() ([]model.Location, error) {
	stmt := `
		SELECT
			location.id,
			location.kind,
			location.name,
			location.info
		FROM location`

	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	rows, err := s.db.Query(ctx, stmt)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var locations []model.Location
	for rows.Next() {
		var location model.Location
		dest := []any{
			&location.ID,
			&location.Kind,
			&location.Name,
			&location.Info,
		}

		err := database.Scan(rows, dest...)
		if err != nil {
			if errors.Is(err, database.ErrRetry) {
				return s.ReadAll()
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
