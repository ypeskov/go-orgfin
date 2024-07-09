package main

import (
	"database/sql"
	"fmt"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/sqlite3"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"ypeskov/go-orgfin/internal/config"
	"ypeskov/go-orgfin/internal/logger"
	"ypeskov/go-orgfin/internal/server"

	_ "github.com/golang-migrate/migrate/v4/database/sqlite3"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

func main() {
	cfg, err := config.New()
	if err != nil {
		fmt.Printf(fmt.Sprintf("cannot read config: %s", err))
		panic(err)
	}

	appLogger := logger.New(cfg)

	err = makeMigration(appLogger, cfg)
	if err != nil {
		panic(fmt.Sprintf("cannot make migration: %s", err))
	}

	appServer := server.New(cfg, appLogger)

	openBrowser(fmt.Sprintf("http://localhost:%s", cfg.Port))

	err = appServer.ListenAndServe()
	if err != nil {
		panic(fmt.Sprintf("cannot start server: %s", err))
	}

}

func openBrowser(url string) {
	var cmd string
	var args []string

	switch runtime.GOOS {
	case "windows":
		cmd = "rundll32"
		args = append(args, "url.dll,FileProtocolHandler", url)
	case "darwin":
		cmd = "open"
		args = append(args, url)
	case "linux":
		cmd = "xdg-open"
		args = append(args, url)
	default:
		fmt.Printf("Unsupported platform: %s\n", runtime.GOOS)
		return
	}
	fmt.Printf("Opening browser at %s\n", url)
	err := exec.Command(cmd, args...).Start()
	if err != nil {
		fmt.Printf("Failed to open browser: %v\n", err)
	}
}

func makeMigration(log *logger.Logger, cfg *config.Config) error {
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
