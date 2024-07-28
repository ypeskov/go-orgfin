package migrations

import (
	"database/sql"
	"embed"
	"errors"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/sqlite3"
	"github.com/golang-migrate/migrate/v4/source/iofs"
	"os"
	"ypeskov/go-password-manager/internal/config"
	"ypeskov/go-password-manager/internal/logger"

	_ "github.com/golang-migrate/migrate/v4/database/sqlite3"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

//go:embed scripts/*.sql
var embeddedMigrations embed.FS

func MakeMigration(log *logger.Logger, cfg *config.Config) error {
	exists, err := checkDatabaseExists(cfg.DbUrl)
	if err != nil {
		log.Fatalf("Error checking if database exists: %s", err)
		return err
	}
	if exists {
		log.Println("Database already exists. Migration not needed.")
		return nil
	}

	err = createDatabase(cfg.DbUrl, log)
	if err != nil {
		return err
	}

	err = applyMigrations(cfg.DbUrl, log)
	if err != nil {
		return err
	}

	log.Println("Migrations applied successfully")
	return nil
}

func checkDatabaseExists(dbUrl string) (bool, error) {
	_, err := os.Stat(dbUrl)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}

func createDatabase(dbUrl string, log *logger.Logger) error {
	_, err := os.Create(dbUrl)
	if err != nil {
		log.Fatalf("Cannot create database file: %s", err)
		return err
	}
	log.Println("Database file created")
	return nil
}

func applyMigrations(dbUrl string, log *logger.Logger) error {
	db, err := sql.Open("sqlite3", dbUrl)
	if err != nil {
		log.Fatalf("Cannot open database: %s", err)
		return err
	}
	defer func() { _ = db.Close() }()

	driver, err := sqlite3.WithInstance(db, &sqlite3.Config{})
	if err != nil {
		log.Fatalf("Could not start SQLite driver: %s", err)
		return err
	}

	sourceDriver, err := iofs.New(embeddedMigrations, "scripts")
	if err != nil {
		log.Fatalf("Could not create iofs driver: %s", err)
		return err
	}

	m, err := migrate.NewWithInstance("iofs", sourceDriver, "sqlite3", driver)
	if err != nil {
		log.Fatalf("Could not create migrate instance: %s", err)
		return err
	}

	if err := m.Up(); err != nil && !errors.Is(err, migrate.ErrNoChange) {
		log.Fatalf("Could not apply migrations: %s", err)
		return err
	}

	return nil
}
