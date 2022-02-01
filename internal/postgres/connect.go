package postgres

import (
	"context"

	"github.com/jackc/pgx/v4/pgxpool"
)

func Connect(uri string) (*pgxpool.Pool, error) {
	// open a database connection pool
	conn, err := pgxpool.Connect(context.Background(), uri)
	if err != nil {
		return nil, err
	}

	// test connection to ensure all is well
	if err = conn.Ping(context.Background()); err != nil {
		conn.Close()
		return nil, err
	}

	return conn, nil
}
