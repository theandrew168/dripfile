package location

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
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
	id string

	kind string
	info []byte

	createdAt time.Time
	version   int
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
		id: l.id,

		kind: l.kind,
		info: encryptedInfoJSON,

		createdAt: l.createdAt,
		version:   l.version,
	}
	return lr, nil
}

func (repo *PostgresRepository) unmarshal(lr locationRow) (*Location, error) {
	switch lr.kind {
	case KindMemory:
		return repo.unmarshalMemory(lr)
	case KindS3:
		return repo.unmarshalS3(lr)
	}

	return nil, fmt.Errorf("unknown location kind: %s", lr.kind)
}

func (repo *PostgresRepository) unmarshalMemory(lr locationRow) (*Location, error) {
	infoJSON, err := repo.box.Decrypt(lr.info)
	if err != nil {
		return nil, err
	}

	var info fileserver.MemoryInfo
	err = json.Unmarshal(infoJSON, &info)
	if err != nil {
		return nil, err
	}

	l := Location{
		id: lr.id,

		kind:       KindMemory,
		memoryInfo: info,

		createdAt: lr.createdAt,
		version:   lr.version,
	}
	return &l, nil
}

func (repo *PostgresRepository) unmarshalS3(lr locationRow) (*Location, error) {
	infoJSON, err := repo.box.Decrypt(lr.info)
	if err != nil {
		return nil, err
	}

	var info fileserver.S3Info
	err = json.Unmarshal(infoJSON, &info)
	if err != nil {
		return nil, err
	}

	l := Location{
		id: lr.id,

		kind:   KindS3,
		s3Info: info,

		createdAt: lr.createdAt,
		version:   lr.version,
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
		lr.id,
		lr.kind,
		lr.info,
	}

	ctx, cancel := context.WithTimeout(context.Background(), database.Timeout)
	defer cancel()

	return database.Exec(repo.conn, ctx, stmt, args...)
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
	defer rows.Close()

	var ls []*Location
	for rows.Next() {
		var lr locationRow
		dest := []any{
			&lr.id,
			&lr.kind,
			&lr.info,
			&lr.createdAt,
			&lr.version,
		}

		err := database.Scan(rows, dest...)
		if err != nil {
			return nil, err
		}

		l, err := repo.unmarshal(lr)
		if err != nil {
			return nil, err
		}

		ls = append(ls, l)
	}

	if rows.Err() != nil {
		return nil, rows.Err()
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

	var lr locationRow
	dest := []any{
		&lr.id,
		&lr.kind,
		&lr.info,
		&lr.createdAt,
		&lr.version,
	}

	ctx, cancel := context.WithTimeout(context.Background(), database.Timeout)
	defer cancel()

	row := repo.conn.QueryRow(ctx, stmt, id)
	err = database.Scan(row, dest...)
	if err != nil {
		return nil, err
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
		lr.id,
		lr.kind,
		lr.info,
		lr.version,
	}

	ctx, cancel := context.WithTimeout(context.Background(), database.Timeout)
	defer cancel()

	row := repo.conn.QueryRow(ctx, stmt, args...)
	err = database.Scan(row, &l.version)
	if err != nil {
		if errors.Is(err, database.ErrNotExist) {
			return database.ErrConflict
		}
		return err
	}

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
		RETURNING id`

	ctx, cancel := context.WithTimeout(context.Background(), database.Timeout)
	defer cancel()

	var deletedID string
	row := repo.conn.QueryRow(ctx, stmt, id)
	return database.Scan(row, &deletedID)
}
