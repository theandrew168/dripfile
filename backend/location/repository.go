package location

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/theandrew168/dripfile/backend/database"
	"github.com/theandrew168/dripfile/backend/fileserver"
	"github.com/theandrew168/dripfile/backend/secret"
)

// repository interface (other code depends on this)
type Repository interface {
	Create(l *Location) error
	Read(id string) (*Location, error)
	List() ([]*Location, error)
	Delete(id string) error
}

// repository implementation (knows about domain internals)
type repository struct {
	conn database.Conn
	box  *secret.Box
}

func NewRepository(conn database.Conn, box *secret.Box) Repository {
	repo := repository{
		conn: conn,
		box:  box,
	}
	return &repo
}

type locationRow struct {
	id string

	kind string
	info []byte
}

func (repo *repository) marshal(l *Location) (locationRow, error) {
	var info any
	switch l.Kind() {
	case KindMemory:
		info = l.MemoryInfo()
	case KindS3:
		info = l.S3Info()
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
		id: l.ID(),

		kind: l.Kind(),
		info: encryptedInfoJSON,
	}
	return lr, nil
}

func (repo *repository) unmarshal(lr locationRow) (*Location, error) {
	switch lr.kind {
	case KindMemory:
		return repo.unmarshalMemory(lr)
	case KindS3:
		return repo.unmarshalS3(lr)
	}

	return nil, fmt.Errorf("unknown location kind: %s", lr.kind)
}

func (repo *repository) unmarshalMemory(lr locationRow) (*Location, error) {
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
	}
	return &l, nil
}

func (repo *repository) unmarshalS3(lr locationRow) (*Location, error) {
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
	}
	return &l, nil
}

func (repo *repository) Create(l *Location) error {
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

func (repo *repository) Read(id string) (*Location, error) {
	stmt := `
		SELECT
			id,
			kind,
			info
		FROM location
		WHERE id = $1`

	var lr locationRow
	dest := []any{
		&lr.id,
		&lr.kind,
		&lr.info,
	}

	ctx, cancel := context.WithTimeout(context.Background(), database.Timeout)
	defer cancel()

	row := repo.conn.QueryRow(ctx, stmt, id)
	err := database.Scan(row, dest...)
	if err != nil {
		return nil, err
	}

	return repo.unmarshal(lr)
}

func (repo *repository) List() ([]*Location, error) {
	stmt := `
		SELECT
			id,
			kind,
			info
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

func (repo *repository) Delete(id string) error {
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
