package server

import (
	"fmt"
	"net/http"
	"strconv"
	"time"
	"ypeskov/go-orgfin/internal/logger"
	"ypeskov/go-orgfin/internal/routes"
	"ypeskov/go-orgfin/services"

	_ "github.com/joho/godotenv/autoload"

	"ypeskov/go-orgfin/internal/config"
	"ypeskov/go-orgfin/internal/database"
)

type Server struct {
	port int
	Db   database.Service
}

func New(cfg *config.Config, logger *logger.Logger) *http.Server {
	port, _ := strconv.Atoi(cfg.Port)
	db := *database.New(cfg)

	servicesManager := services.NewServiceManager(&db, logger)
	echo := routes.RegisterRoutes(logger, servicesManager)

	server := &http.Server{
		Addr:         fmt.Sprintf(":%d", port),
		Handler:      echo,
		IdleTimeout:  time.Minute,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
	}

	return server
}
