package postgres

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/theandrew168/dripfile/internal/database"
	"github.com/theandrew168/dripfile/internal/location"
	"github.com/theandrew168/dripfile/internal/location/fileserver/s3"
	"github.com/theandrew168/dripfile/internal/secret"
)

type locationModel struct {
	id   string
	kind string
	info []byte
}

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

func (store *Storage) marshal(l *location.Location) (locationModel, error) {
	var info any
	switch l.Kind() {
	case location.KindS3:
		info = l.S3Info()
	}

	infoJSON, err := json.Marshal(info)
	if err != nil {
		return locationModel{}, err
	}

	encryptedInfoJSON, err := store.box.Encrypt(infoJSON)
	if err != nil {
		return locationModel{}, err
	}

	lm := locationModel{
		id:   l.ID(),
		kind: l.Kind(),
		info: encryptedInfoJSON,
	}
	return lm, nil
}

func (store *Storage) unmarshal(lm locationModel) (*location.Location, error) {
	switch lm.kind {
	case location.KindS3:
		return store.unmarshalS3(lm)
	}

	return nil, fmt.Errorf("unknown location kind: %s", lm.kind)
}

func (store *Storage) unmarshalS3(lm locationModel) (*location.Location, error) {
	infoJSON, err := store.box.Decrypt(lm.info)
	if err != nil {
		return nil, err
	}

	var info s3.Info
	err = json.Unmarshal(infoJSON, &info)
	if err != nil {
		return nil, err
	}

	return location.UnmarshalS3FromStorage(lm.id, info)
}

func (store *Storage) Create(l *location.Location) error {
	stmt := `
		INSERT INTO location
			(id, kind, info)
		VALUES
			($1, $2, $3)`

	lm, err := store.marshal(l)
	if err != nil {
		return err
	}

	args := []any{
		lm.id,
		lm.kind,
		lm.info,
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

	var lm locationModel
	dest := []any{
		&lm.id,
		&lm.kind,
		&lm.info,
	}

	ctx, cancel := context.WithTimeout(context.Background(), database.Timeout)
	defer cancel()

	row := store.conn.QueryRow(ctx, stmt, id)
	err := database.Scan(row, dest...)
	if err != nil {
		return nil, err
	}

	return store.unmarshal(lm)
}

func (store *Storage) List() ([]*location.Location, error) {
	stmt := `
		SELECT
			id,
			kind,
			info
		FROM location`

	ctx, cancel := context.WithTimeout(context.Background(), database.Timeout)
	defer cancel()

	rows, err := store.conn.Query(ctx, stmt)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var ls []*location.Location
	for rows.Next() {
		var lm locationModel
		dest := []any{
			&lm.id,
			&lm.kind,
			&lm.info,
		}

		err := database.Scan(rows, dest...)
		if err != nil {
			return nil, err
		}

		l, err := store.unmarshal(lm)
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
