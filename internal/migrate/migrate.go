package migrate

import (
	"context"
	"embed"
	"fmt"
	"io/fs"
	"sort"

	"github.com/theandrew168/dripfile/internal/jsonlog"
	"github.com/theandrew168/dripfile/internal/postgresql"
)

//go:embed migration
var migrationFS embed.FS

func Migrate(logger *jsonlog.Logger, db postgresql.Conn) error {
	ctx := context.Background()

	// attempt to create extensions and ignore errors
	// (local container needs it, deployed version gets it via ansible)
	exts := []string{
		"citext",
		"pgcrypto",
	}
	for _, ext := range exts {
		stmt := fmt.Sprintf("CREATE EXTENSION IF NOT EXISTS %s", ext)
		_, err := db.Exec(ctx, stmt)
		if err != nil {
			return err
		}
	}

	// create migrations table if it doesn't exist
	_, err := db.Exec(ctx, `
		CREATE TABLE IF NOT EXISTS migration (
			id SERIAL PRIMARY KEY,
			name TEXT NOT NULL UNIQUE
		)`)
	if err != nil {
		return err
	}

	// get migrations that are already applied
	rows, err := db.Query(ctx, "SELECT name FROM migration")
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
	subdir, _ := fs.Sub(migrationFS, "migration")
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
	for _, name := range missing {
		logger.Info("applying migration", map[string]string{
			"name": name,
		})

		// apply the missing ones
		sql, err := fs.ReadFile(subdir, name)
		if err != nil {
			return err
		}

		// apply each migration in a transaction
		tx, err := db.Begin(context.Background())
		if err != nil {
			return err
		}
		defer tx.Rollback(context.Background())

		_, err = tx.Exec(ctx, string(sql))
		if err != nil {
			return err
		}

		// update migrations table
		_, err = tx.Exec(ctx, "INSERT INTO migration (name) VALUES ($1)", name)
		if err != nil {
			return err
		}

		err = tx.Commit(context.Background())
		if err != nil {
			return err
		}
	}

	logger.Info("migrations up to date", nil)
	return nil
}
