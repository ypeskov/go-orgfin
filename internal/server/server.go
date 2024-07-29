package server

import (
	"fmt"
	"net/http"
	"strconv"
	"time"
	"ypeskov/go-password-manager/internal/config"
	"ypeskov/go-password-manager/internal/database"
	"ypeskov/go-password-manager/internal/logger"
	"ypeskov/go-password-manager/internal/routes"
	"ypeskov/go-password-manager/services"

	_ "github.com/joho/godotenv/autoload"
)

type Server struct {
	port int
	Db   database.DbService
}

func New(cfg *config.Config, logger *logger.Logger) *http.Server {
	port, _ := strconv.Atoi(cfg.Port)
	db := *database.New(cfg)

	servicesManager := services.NewServiceManager(&db, logger)
	echo := routes.RegisterRoutes(logger, servicesManager, cfg)

	server := &http.Server{
		Addr:         fmt.Sprintf(":%d", port),
		Handler:      echo,
		IdleTimeout:  time.Minute,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
	}

	return server
}
