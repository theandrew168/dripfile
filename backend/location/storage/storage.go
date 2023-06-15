package storage

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/theandrew168/dripfile/backend/database"
	"github.com/theandrew168/dripfile/backend/location"
	"github.com/theandrew168/dripfile/backend/location/fileserver/memory"
	"github.com/theandrew168/dripfile/backend/location/fileserver/s3"
	"github.com/theandrew168/dripfile/backend/secret"
)

type Storage struct {
	conn database.Conn
	box  *secret.Box
}

func New(conn database.Conn, box *secret.Box) *Storage {
	store := Storage{
		conn: conn,
		box:  box,
	}
	return &store
}

type locationRow struct {
	id string

	kind string
	info []byte
}

func (store *Storage) marshal(l *location.Location) (locationRow, error) {
	var info any
	switch l.Kind() {
	case location.KindMemory:
		info = l.MemoryInfo()
	case location.KindS3:
		info = l.S3Info()
	}

	infoJSON, err := json.Marshal(info)
	if err != nil {
		return locationRow{}, err
	}

	encryptedInfoJSON, err := store.box.Encrypt(infoJSON)
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

func (store *Storage) unmarshal(lr locationRow) (*location.Location, error) {
	switch lr.kind {
	case location.KindMemory:
		return store.unmarshalMemory(lr)
	case location.KindS3:
		return store.unmarshalS3(lr)
	}

	return nil, fmt.Errorf("unknown location kind: %s", lr.kind)
}

func (store *Storage) unmarshalMemory(lr locationRow) (*location.Location, error) {
	infoJSON, err := store.box.Decrypt(lr.info)
	if err != nil {
		return nil, err
	}

	var info memory.Info
	err = json.Unmarshal(infoJSON, &info)
	if err != nil {
		return nil, err
	}

	return location.UnmarshalMemoryFromStorage(lr.id, info)
}

func (store *Storage) unmarshalS3(lr locationRow) (*location.Location, error) {
	infoJSON, err := store.box.Decrypt(lr.info)
	if err != nil {
		return nil, err
	}

	var info s3.Info
	err = json.Unmarshal(infoJSON, &info)
	if err != nil {
		return nil, err
	}

	return location.UnmarshalS3FromStorage(lr.id, info)
}

func (store *Storage) Create(l *location.Location) error {
	stmt := `
		INSERT INTO location
			(id, kind, info)
		VALUES
			($1, $2, $3)`

	lr, err := store.marshal(l)
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

	return database.Exec(store.conn, ctx, stmt, args...)
}

func (store *Storage) Read(id string) (*location.Location, error) {
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

	row := store.conn.QueryRow(ctx, stmt, id)
	err := database.Scan(row, dest...)
	if err != nil {
		return nil, err
	}

	return store.unmarshal(lr)
}

func (store *Storage) List() ([]*location.Location, error) {
	stmt := `
		SELECT
			id,
			kind,
			info
		FROM location
		ORDER BY created_at ASC`

	ctx, cancel := context.WithTimeout(context.Background(), database.Timeout)
	defer cancel()

	rows, err := store.conn.Query(ctx, stmt)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var ls []*location.Location
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

		l, err := store.unmarshal(lr)
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

func (store *Storage) Delete(id string) error {
	stmt := `
		DELETE FROM location
		WHERE id = $1
		RETURNING id`

	ctx, cancel := context.WithTimeout(context.Background(), database.Timeout)
	defer cancel()

	var deletedID string
	row := store.conn.QueryRow(ctx, stmt, id)
	return database.Scan(row, &deletedID)
}
