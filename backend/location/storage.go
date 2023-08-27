package location

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/theandrew168/dripfile/backend/database"
	"github.com/theandrew168/dripfile/backend/fileserver"
	"github.com/theandrew168/dripfile/backend/secret"
)

// storage interface (other code depends on this)
type Storage interface {
	Create(l *Location) error
	Read(id string) (*Location, error)
	List() ([]*Location, error)
	Delete(id string) error
}

// storage implementation (knows about domain internals)
type storage struct {
	conn database.Conn
	box  *secret.Box
}

func NewStorage(conn database.Conn, box *secret.Box) Storage {
	store := storage{
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

func (store *storage) marshal(l *Location) (locationRow, error) {
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

func (store *storage) unmarshal(lr locationRow) (*Location, error) {
	switch lr.kind {
	case KindMemory:
		return store.unmarshalMemory(lr)
	case KindS3:
		return store.unmarshalS3(lr)
	}

	return nil, fmt.Errorf("unknown location kind: %s", lr.kind)
}

func (store *storage) unmarshalMemory(lr locationRow) (*Location, error) {
	infoJSON, err := store.box.Decrypt(lr.info)
	if err != nil {
		return nil, err
	}

	var info fileserver.MemoryInfo
	err = json.Unmarshal(infoJSON, &info)
	if err != nil {
		return nil, err
	}

	return UnmarshalMemoryFromStorage(lr.id, info)
}

func (store *storage) unmarshalS3(lr locationRow) (*Location, error) {
	infoJSON, err := store.box.Decrypt(lr.info)
	if err != nil {
		return nil, err
	}

	var info fileserver.S3Info
	err = json.Unmarshal(infoJSON, &info)
	if err != nil {
		return nil, err
	}

	return UnmarshalS3FromStorage(lr.id, info)
}

func (store *storage) Create(l *Location) error {
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

func (store *storage) Read(id string) (*Location, error) {
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

func (store *storage) List() ([]*Location, error) {
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

func (store *storage) Delete(id string) error {
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
