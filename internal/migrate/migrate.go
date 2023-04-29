package migrate

import (
	"context"
	"embed"
	"fmt"
	"io/fs"
	"sort"

	"github.com/theandrew168/dripfile/internal/database"
)

func Migrate(conn database.Conn, files embed.FS) error {
	ctx := context.Background()

	// attempt to create extensions (requires superuser privileges)
	// (works against local container, deployed DB needs these at setup)
	exts := []string{
		"citext",
		"pgcrypto",
	}
	for _, ext := range exts {
		stmt := fmt.Sprintf("CREATE EXTENSION IF NOT EXISTS %s", ext)
		_, err := conn.Exec(ctx, stmt)
		if err != nil {
			return err
		}
	}

	// create migration table if it doesn't exist
	_, err := conn.Exec(ctx, `
		CREATE TABLE IF NOT EXISTS migration (
			id SERIAL PRIMARY KEY,
			name TEXT NOT NULL UNIQUE
		)`)
	if err != nil {
		return err
	}

	// get migrations that are already applied
	rows, err := conn.Query(ctx, "SELECT name FROM migration")
	if err != nil {
		return err
	}
	defer rows.Close()

	applied := make(map[string]bool)
	for rows.Next() {
		var name string
		err = rows.Scan(&name)
		if err != nil {
			return err
		}
		applied[name] = true
	}

	// get migrations that should be applied (from migrations FS)
	subdir, _ := fs.Sub(files, "migration")
	migrations, err := fs.ReadDir(subdir, ".")
	if err != nil {
		return err
	}

	// determine missing migrations
	var missing []string
	for _, migration := range migrations {
		name := migration.Name()
		if _, ok := applied[name]; !ok {
			missing = append(missing, name)
		}
	}

	// sort missing migrations to preserve order
	sort.Strings(missing)

	// apply each missing migration
	for _, name := range missing {
		sql, err := fs.ReadFile(subdir, name)
		if err != nil {
			return err
		}

		// apply each migration in a transaction
		tx, err := conn.Begin(context.Background())
		if err != nil {
			return err
		}
		defer tx.Rollback(context.Background())

		_, err = tx.Exec(ctx, string(sql))
		if err != nil {
			return err
		}

		// update migration table
		_, err = tx.Exec(ctx, "INSERT INTO migration (name) VALUES ($1)", name)
		if err != nil {
			return err
		}

		err = tx.Commit(context.Background())
		if err != nil {
			return err
		}
	}

	return nil
}
