package main

import (
	"fmt"
	"ypeskov/go-orgfin/internal/config"
	"ypeskov/go-orgfin/internal/logger"
	"ypeskov/go-orgfin/internal/server"
)

func main() {
	cfg, err := config.New()
	if err != nil {
		fmt.Sprintf(fmt.Sprintf("cannot read config: %s", err))
		panic(err)
	}

	appLogger := logger.New(cfg)

	appServer := server.New(cfg, appLogger)
	err = appServer.ListenAndServe()
	if err != nil {
		panic(fmt.Sprintf("cannot start server: %s", err))
	}
}
