package main

import (
	"context"
	"embed"
	"flag"
	"io/fs"
	"log"
	"os"
	"sort"

	"github.com/jackc/pgx/v4/pgxpool"

	"github.com/theandrew168/dripfile/internal/config"
)

//go:embed migrations
var migrationsFS embed.FS

func main() {
	logger := log.New(os.Stdout, "", log.Lshortfile)

	// check for config file flag
	conf := flag.String("conf", "dripfile.conf", "app config file")
	flag.Parse()

	// load user-defined config (if specified), else use defaults
	cfg, err := config.ReadFile(*conf)
	if err != nil {
		logger.Fatalln(err)
	}

	// open a database connection pool
	conn, err := pgxpool.Connect(context.Background(), cfg.DatabaseURI)
	if err != nil {
		logger.Fatalln(err)
	}
	defer conn.Close()

	// test connection to ensure all is well
	if err = conn.Ping(context.Background()); err != nil {
		logger.Fatalln(err)
	}

	// apply database migrations
	migrations, _ := fs.Sub(migrationsFS, "migrations")
	if err = compareAndApplyMigrations(conn, migrations, logger); err != nil {
		logger.Fatalln(err)
	}
}

func compareAndApplyMigrations(conn *pgxpool.Pool, migrationsFS fs.FS, logger *log.Logger) error {
	ctx := context.Background()

	// create migrations table if it doesn't exist
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

	// get migrations that should be applied (from migrations/ dir)
	migrations, err := fs.ReadDir(migrationsFS, ".")
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
		logger.Printf("applying: %s\n", name)

		// apply the missing ones
		sql, err := fs.ReadFile(migrationsFS, name)
		if err != nil {
			return err
		}
		_, err = conn.Exec(ctx, string(sql))
		if err != nil {
			return err
		}

		// update migrations table
		_, err = conn.Exec(ctx, "INSERT INTO migration (name) VALUES ($1)", name)
		if err != nil {
			return err
		}
	}

	logger.Println("migrations up to date")
	return nil
}
