package database

import (
	"fmt"
	"github.com/jmoiron/sqlx"
	log "github.com/sirupsen/logrus"
	"ypeskov/go-orgfin/internal/config"

	_ "github.com/joho/godotenv/autoload"
	_ "github.com/mattn/go-sqlite3"
)

// Service represents a service that interacts with a database.
type Service interface {
	// Close terminates the database connection.
	// It returns an error if the connection cannot be closed.
	Close() error
}

type service struct {
	db    *sqlx.DB
	dbUrl string
}

var (
	dbInstance *service
)

func New(cfg *config.Config) Service {
	// Reuse Connection
	if dbInstance != nil {
		return dbInstance
	}

	//db, err := sqlx.Connect("postgres", dbConnStr)
	db, err := sqlx.Open("sqlite3", cfg.DbUrl)
	if err != nil {
		log.Fatal(err)
	}
	dbInstance = &service{
		db:    db,
		dbUrl: cfg.DbUrl,
	}
	log.Info(fmt.Sprintf("Connected to database: %#v\n", dbInstance))

	return dbInstance
}

// Close closes the database connection.
// It logs a message indicating the disconnection from the specific database.
// If the connection is successfully closed, it returns nil.
// If an error occurs while closing the connection, it returns the error.
func (s *service) Close() error {
	log.Printf("Disconnected from database: %s", s.dbUrl)
	return s.db.Close()
}
