package migrate

import (
	"context"
	"embed"
	"fmt"
	"io/fs"
	"sort"

	"github.com/theandrew168/dripfile/pkg/jsonlog"
	"github.com/theandrew168/dripfile/pkg/postgres"
)

//go:embed migration
var migrationFS embed.FS

func Migrate(pg postgres.Interface, logger *jsonlog.Logger) error {
	ctx := context.Background()

	// attempt to create extensions and ignore errors
	// (local container needs it, deployed version gets it via ansible)
	exts := []string{
		"citext",
		"pgcrypto",
	}
	for _, ext := range exts {
		stmt := fmt.Sprintf("CREATE EXTENSION IF NOT EXISTS %s", ext)
		_, err := pg.Exec(ctx, stmt)
		if err != nil {
			return err
		}
	}

	// create migrations table if it doesn't exist
	_, err := pg.Exec(ctx, `
		CREATE TABLE IF NOT EXISTS migration (
			id SERIAL PRIMARY KEY,
			name TEXT NOT NULL UNIQUE
		)`)
	if err != nil {
		return err
	}

	// get migrations that are already applied
	rows, err := pg.Query(ctx, "SELECT name FROM migration")
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
		logger.PrintInfo("applying migration", map[string]string{
			"name": name,
		})

		// apply the missing ones
		sql, err := fs.ReadFile(subdir, name)
		if err != nil {
			return err
		}
		_, err = pg.Exec(ctx, string(sql))
		if err != nil {
			return err
		}

		// update migrations table
		_, err = pg.Exec(ctx, "INSERT INTO migration (name) VALUES ($1)", name)
		if err != nil {
			return err
		}
	}

	logger.PrintInfo("migrations up to date", nil)
	return nil
}
