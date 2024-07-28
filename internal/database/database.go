package database

import (
	"fmt"
	"github.com/jmoiron/sqlx"
	log "github.com/sirupsen/logrus"
	"ypeskov/go-password-manager/internal/config"

	_ "github.com/joho/godotenv/autoload"
	_ "github.com/mattn/go-sqlite3"
)

type DbService struct {
	Db    *sqlx.DB
	DbUrl string
}

var (
	dbInstance *DbService
)

func New(cfg *config.Config) *DbService {
	// Reuse Connection
	if dbInstance != nil {
		return dbInstance
	}

	db, err := sqlx.Open("sqlite3", cfg.DbUrl)
	if err != nil {
		log.Fatal(err)
	}
	dbInstance = &DbService{
		Db:    db,
		DbUrl: cfg.DbUrl,
	}
	log.Info(fmt.Sprintf("Connected to database: %#v\n", dbInstance))

	return dbInstance
}
