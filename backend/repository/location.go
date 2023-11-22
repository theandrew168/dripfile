package repository

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"

	"github.com/theandrew168/dripfile/backend/database"
	"github.com/theandrew168/dripfile/backend/domain"
	"github.com/theandrew168/dripfile/backend/fileserver"
	"github.com/theandrew168/dripfile/backend/secret"
)

// ensure LocationRepository interface is satisfied
var _ LocationRepository = (*PostgresLocationRepository)(nil)

type LocationRepository interface {
	Create(location *domain.Location) error
	List() ([]*domain.Location, error)
	Read(id uuid.UUID) (*domain.Location, error)
	Delete(location *domain.Location) error
}

type Location struct {
	ID uuid.UUID `db:"id"`

	Kind domain.LocationKind `db:"kind"`
	Info []byte              `db:"info"`

	CreatedAt time.Time `db:"created_at"`
	Version   int       `db:"version"`
}

type PostgresLocationRepository struct {
	conn database.Conn
	box  *secret.Box
}

func NewPostgresLocationRepository(conn database.Conn, box *secret.Box) *PostgresLocationRepository {
	repo := PostgresLocationRepository{
		conn: conn,
		box:  box,
	}
	return &repo
}

func (repo *PostgresLocationRepository) marshal(location *domain.Location) (Location, error) {
	info := location.Info()
	infoJSON, err := json.Marshal(info)
	if err != nil {
		return Location{}, err
	}

	encryptedInfoJSON, err := repo.box.Encrypt(infoJSON)
	if err != nil {
		return Location{}, err
	}

	row := Location{
		ID: location.ID(),

		Kind: location.Kind(),
		Info: encryptedInfoJSON,

		CreatedAt: location.CreatedAt(),
		Version:   location.Version(),
	}
	return row, nil
}

func (repo *PostgresLocationRepository) unmarshal(row Location) (*domain.Location, error) {
	switch row.Kind {
	case domain.LocationKindMemory:
		return repo.unmarshalMemory(row)
	case domain.LocationKindS3:
		return repo.unmarshalS3(row)
	}

	return nil, fmt.Errorf("unknown location kind: %s", row.Kind)
}

func (repo *PostgresLocationRepository) unmarshalMemory(row Location) (*domain.Location, error) {
	infoJSON, err := repo.box.Decrypt(row.Info)
	if err != nil {
		return nil, err
	}

	var info fileserver.MemoryInfo
	err = json.Unmarshal(infoJSON, &info)
	if err != nil {
		return nil, err
	}

	location := domain.LoadMemoryLocation(row.ID, info, row.CreatedAt, row.Version)
	return location, nil
}

func (repo *PostgresLocationRepository) unmarshalS3(row Location) (*domain.Location, error) {
	infoJSON, err := repo.box.Decrypt(row.Info)
	if err != nil {
		return nil, err
	}

	var info fileserver.S3Info
	err = json.Unmarshal(infoJSON, &info)
	if err != nil {
		return nil, err
	}

	location := domain.LoadS3Location(row.ID, info, row.CreatedAt, row.Version)
	return location, nil
}

func (repo *PostgresLocationRepository) Create(location *domain.Location) error {
	stmt := `
		INSERT INTO location
			(id, kind, info)
		VALUES
			($1, $2, $3)`

	row, err := repo.marshal(location)
	if err != nil {
		return err
	}

	args := []any{
		row.ID,
		row.Kind,
		row.Info,
	}

	ctx, cancel := context.WithTimeout(context.Background(), database.Timeout)
	defer cancel()

	_, err = repo.conn.Exec(ctx, stmt, args...)
	if err != nil {
		return checkCreateError(err)
	}

	// TODO: set createdAt and version fields
	return nil
}

func (repo *PostgresLocationRepository) List() ([]*domain.Location, error) {
	stmt := `
		SELECT
			id,
			kind,
			info,
			created_at,
			version
		FROM location
		ORDER BY created_at ASC`

	ctx, cancel := context.WithTimeout(context.Background(), database.Timeout)
	defer cancel()

	rows, err := repo.conn.Query(ctx, stmt)
	if err != nil {
		return nil, err
	}

	locationRows, err := pgx.CollectRows(rows, pgx.RowToStructByName[Location])
	if err != nil {
		return nil, checkListError(err)
	}

	var locations []*domain.Location
	for _, row := range locationRows {
		location, err := repo.unmarshal(row)
		if err != nil {
			return nil, err
		}

		locations = append(locations, location)
	}

	return locations, nil
}

func (repo *PostgresLocationRepository) Read(id uuid.UUID) (*domain.Location, error) {
	stmt := `
		SELECT
			id,
			kind,
			info,
			created_at,
			version
		FROM location
		WHERE id = $1`

	ctx, cancel := context.WithTimeout(context.Background(), database.Timeout)
	defer cancel()

	rows, err := repo.conn.Query(ctx, stmt, id)
	if err != nil {
		return nil, err
	}

	row, err := pgx.CollectOneRow(rows, pgx.RowToStructByName[Location])
	if err != nil {
		return nil, checkReadError(err)
	}

	return repo.unmarshal(row)
}

func (repo *PostgresLocationRepository) Delete(location *domain.Location) error {
	stmt := `
		DELETE FROM location
		WHERE id = $1
		RETURNING version`

	err := location.CheckDelete()
	if err != nil {
		return err
	}

	ctx, cancel := context.WithTimeout(context.Background(), database.Timeout)
	defer cancel()

	rows, err := repo.conn.Query(ctx, stmt, location.ID())
	if err != nil {
		return err
	}

	_, err = pgx.CollectOneRow(rows, pgx.RowTo[int])
	if err != nil {
		return checkDeleteError(err)
	}

	return nil
}
