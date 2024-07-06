package main

import (
	"fmt"
	log "github.com/sirupsen/logrus"
	"ypeskov/go-orgfin/internal/config"
	"ypeskov/go-orgfin/internal/server"
)

func main() {
	cfg, err := config.New()
	if err != nil {
		log.Panic(fmt.Sprintf("cannot read config: %s", err))
	}

	appServer := server.New(cfg)
	err = appServer.ListenAndServe()
	if err != nil {
		panic(fmt.Sprintf("cannot start server: %s", err))
	}
}
