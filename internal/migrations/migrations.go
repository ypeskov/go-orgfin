package migrations

import (
	"database/sql"
	"fmt"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/sqlite3"
	"os"
	"path/filepath"
	"ypeskov/go-orgfin/internal/config"
	"ypeskov/go-orgfin/internal/logger"

	_ "github.com/golang-migrate/migrate/v4/database/sqlite3"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

func MakeMigration(log *logger.Logger, cfg *config.Config) error {
	if _, err := os.Stat(cfg.DbUrl); os.IsNotExist(err) {
		file, err := os.Create(cfg.DbUrl)
		if err != nil {
			log.Fatalf("cannot create database file: %s", err)
			return err
		}
		file.Close()
		log.Println("Database file created")
	}

	db, err := sql.Open("sqlite3", cfg.DbUrl)
	if err != nil {
		log.Fatalf("cannot open database: %s", err)
		return err
	}
	defer db.Close()

	driver, err := sqlite3.WithInstance(db, &sqlite3.Config{})
	if err != nil {
		log.Fatalf("could not start SQLite driver: %s", err)
		return err
	}

	migrationsPath, err := filepath.Abs("migrations")
	log.Info(fmt.Sprintf("Migrations path: %s", migrationsPath))
	if err != nil {
		log.Fatalf("could not get absolute path to migrations directory: %s", err)
		return err
	}

	m, err := migrate.NewWithDatabaseInstance(
		"file://"+migrationsPath,
		"sqlite3", driver)
	if err != nil {
		log.Fatalf("could not create migrate instance: %s", err)
		return err
	}

	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		log.Fatalf("could not apply migrations: %s", err)
		return err
	}

	log.Println("Migrations applied successfully")

	return nil
}
