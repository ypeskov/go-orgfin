package server

import (
	"fmt"
	"net/http"
	"strconv"
	"time"

	_ "github.com/joho/godotenv/autoload"

	"ypeskov/go-orgfin/internal/config"
	"ypeskov/go-orgfin/internal/database"
)

type Server struct {
	port int

	db database.Service
}

func New(cfg *config.Config) *http.Server {
	port, _ := strconv.Atoi(cfg.Port)
	NewServer := &Server{
		port: port,

		db: database.New(cfg),
	}

	// Declare Server config
	server := &http.Server{
		Addr:         fmt.Sprintf(":%d", NewServer.port),
		Handler:      NewServer.RegisterRoutes(),
		IdleTimeout:  time.Minute,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
	}

	return server
}
