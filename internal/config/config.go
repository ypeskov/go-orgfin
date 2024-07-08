package config

import (
	"fmt"
	"github.com/caarlos0/env/v10"
	"github.com/joho/godotenv"
	log "github.com/sirupsen/logrus"
)

type Config struct {
	Port                        string `env:"PORT" envDefault:"3000"`
	LogLevel                    string `env:"LOG_LEVEL" envDefault:"INFO"`
	DbUrl                       string `env:"DATABASE_URL" envDefault:"db.sqlite3"`
	SecretKey                   string `env:"SECRET_KEY" envDefault:"secret"`
	AccessTokenLifetimeMinutes  int    `env:"ACCESS_TOKEN_LIFETIME_MINUTES" envDefault:"5"`
	RefreshTokenLifetimeMinutes int    `env:"REFRESH_TOKEN_LIFETIME_MINUTES" envDefault:"1440"`
}

func New() (*Config, error) {
	_ = godotenv.Load(".env")

	cfg := &Config{}
	if err := env.Parse(cfg); err != nil {
		log.Panic("Error parsing env vars: %v", err)

		return nil, err
	}
	fmt.Printf("Config: %+v\n", cfg)

	return cfg, nil
}
