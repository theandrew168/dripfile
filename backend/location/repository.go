package location

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/theandrew168/dripfile/backend/database"
	"github.com/theandrew168/dripfile/backend/fileserver"
	"github.com/theandrew168/dripfile/backend/secret"
)

// ensure Repository interface is satisfied
var _ Repository = (*PostgresRepository)(nil)

// repository interface (other code depends on this)
type Repository interface {
	Create(l *Location) error
	List() ([]*Location, error)
	Read(id string) (*Location, error)
	Update(l *Location) error
	Delete(id string) error
}

// repository implementation (knows about domain internals)
type PostgresRepository struct {
	conn database.Conn
	box  *secret.Box
}

func NewRepository(conn database.Conn, box *secret.Box) *PostgresRepository {
	repo := PostgresRepository{
		conn: conn,
		box:  box,
	}
	return &repo
}

type locationRow struct {
	ID string `db:"id"`

	Kind string `db:"kind"`
	Info []byte `db:"info"`

	CreatedAt time.Time `db:"created_at"`
	Version   int       `db:"version"`
}

func (repo *PostgresRepository) marshal(l *Location) (locationRow, error) {
	var info any
	switch l.Kind() {
	case KindMemory:
		info = l.memoryInfo
	case KindS3:
		info = l.s3Info
	}

	infoJSON, err := json.Marshal(info)
	if err != nil {
		return locationRow{}, err
	}

	encryptedInfoJSON, err := repo.box.Encrypt(infoJSON)
	if err != nil {
		return locationRow{}, err
	}

	lr := locationRow{
		ID: l.id,

		Kind: l.kind,
		Info: encryptedInfoJSON,

		CreatedAt: l.createdAt,
		Version:   l.version,
	}
	return lr, nil
}

func (repo *PostgresRepository) unmarshal(lr locationRow) (*Location, error) {
	switch lr.Kind {
	case KindMemory:
		return repo.unmarshalMemory(lr)
	case KindS3:
		return repo.unmarshalS3(lr)
	}

	return nil, fmt.Errorf("unknown location kind: %s", lr.Kind)
}

func (repo *PostgresRepository) unmarshalMemory(lr locationRow) (*Location, error) {
	infoJSON, err := repo.box.Decrypt(lr.Info)
	if err != nil {
		return nil, err
	}

	var info fileserver.MemoryInfo
	err = json.Unmarshal(infoJSON, &info)
	if err != nil {
		return nil, err
	}

	l := Location{
		id: lr.ID,

		kind:       KindMemory,
		memoryInfo: info,

		createdAt: lr.CreatedAt,
		version:   lr.Version,
	}
	return &l, nil
}

func (repo *PostgresRepository) unmarshalS3(lr locationRow) (*Location, error) {
	infoJSON, err := repo.box.Decrypt(lr.Info)
	if err != nil {
		return nil, err
	}

	var info fileserver.S3Info
	err = json.Unmarshal(infoJSON, &info)
	if err != nil {
		return nil, err
	}

	l := Location{
		id: lr.ID,

		kind:   KindS3,
		s3Info: info,

		createdAt: lr.CreatedAt,
		version:   lr.Version,
	}
	return &l, nil
}

func (repo *PostgresRepository) Create(l *Location) error {
	stmt := `
		INSERT INTO location
			(id, kind, info)
		VALUES
			($1, $2, $3)`

	lr, err := repo.marshal(l)
	if err != nil {
		return err
	}

	args := []any{
		lr.ID,
		lr.Kind,
		lr.Info,
	}

	ctx, cancel := context.WithTimeout(context.Background(), database.Timeout)
	defer cancel()

	_, err = repo.conn.Exec(ctx, stmt, args...)
	if err != nil {
		return database.CheckCreateError(err)
	}

	return nil
}

func (repo *PostgresRepository) List() ([]*Location, error) {
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

	lrs, err := pgx.CollectRows(rows, pgx.RowToStructByName[locationRow])
	if err != nil {
		return nil, err
	}

	var ls []*Location
	for _, lr := range lrs {
		l, err := repo.unmarshal(lr)
		if err != nil {
			return nil, err
		}

		ls = append(ls, l)
	}

	return ls, nil
}

func (repo *PostgresRepository) Read(id string) (*Location, error) {
	_, err := uuid.Parse(id)
	if err != nil {
		return nil, database.ErrInvalidUUID
	}

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

	lr, err := pgx.CollectOneRow(rows, pgx.RowToStructByName[locationRow])
	if err != nil {
		return nil, database.CheckReadError(err)
	}

	return repo.unmarshal(lr)
}

func (repo *PostgresRepository) Update(l *Location) error {
	stmt := `
		UPDATE location
		SET
			kind = $2,
			info = $3,
			version = version + 1
		WHERE id = $1
		  AND version = $4
		RETURNING version`

	lr, err := repo.marshal(l)
	if err != nil {
		return err
	}

	args := []any{
		lr.ID,
		lr.Kind,
		lr.Info,
		lr.Version,
	}

	ctx, cancel := context.WithTimeout(context.Background(), database.Timeout)
	defer cancel()

	rows, err := repo.conn.Query(ctx, stmt, args...)
	if err != nil {
		return err
	}

	version, err := pgx.CollectOneRow(rows, pgx.RowTo[int])
	if err != nil {
		return database.CheckUpdateError(err)
	}

	l.version = version
	return err
}

func (repo *PostgresRepository) Delete(id string) error {
	_, err := uuid.Parse(id)
	if err != nil {
		return database.ErrInvalidUUID
	}

	stmt := `
		DELETE FROM location
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
