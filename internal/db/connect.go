package db

import (
	"context"
	"database/sql"
	"fmt"
	"log/slog"
	"path/filepath"

	"github.com/ncruces/go-sqlite3"
	"github.com/ncruces/go-sqlite3/driver"
	_ "github.com/ncruces/go-sqlite3/embed"

	"github.com/pressly/goose/v3"
)

func Connect(ctx context.Context, data_dir string) (*sql.DB, error) {
	if data_dir == "" {
		return nil, fmt.Errorf("data.dir is not set")
	}

	db_path := filepath.Join(data_dir, "wm.db")

	// Define pragmas, para melhor desempenho
	pragmas := []string{
		"PRAGMA foreign_keys = ON;",
		"PRAGMA journal_mode = WAL;",
		"PRAGMA page_size = 4096;",
		"PRAGMA cache_size = -8000;",
		"PRAGMA synchronous = NORMAL;",
		"PRAGMA secure_delete = ON;",
	}
	db, err := driver.Open(db_path, func(c *sqlite3.Conn) error {
		for _, pragma := range pragmas {
			if err := c.Exec(pragma); err != nil {
				return fmt.Errorf("failed to set pragma `%s`: %w", pragma, err)
			}
		}
		return nil
	})
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}
	// Verifica conexão
	if err = db.PingContext(ctx); err != nil {
		db.Close()
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	goose.SetBaseFS(FS)

	if err := goose.SetDialect("sqlite3"); err != nil {
		slog.Error("Failed to set dialect", "error", err)
		return nil, fmt.Errorf("failed to set dialect: %w", err)
	}

	if err := goose.Up(db, "migrations"); err != nil {
		slog.Error("Failed to apply migrations", "error", err)
		return nil, fmt.Errorf("failed to apply migrations: %w", err)
	}
	return db, nil
}
