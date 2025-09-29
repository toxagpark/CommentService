package config

import (
	"errors"
	"fmt"

	"github.com/caarlos0/env/v11"
)

var (
	ErrConfig = errors.New("config error")
)

type (
	Config struct {
		Postgres Postgres
	}

	Postgres struct {
		POSTGRES_HOST     string `env:"POSTGRES_HOST" envDefault:"localhost"`
		POSTGRES_PORT     string `env:"POSTGRES_PORT" envDefault:"5433"`
		POSTGRES_DB       string `env:"POSTGRES_DB,required"`
		POSTGRES_USER     string `env:"POSTGRES_USER,required"`
		POSTGRES_PASSWORD string `env:"POSTGRES_PASSWORD,required"`
		POSTGRES_SSLMODE  string `env:"POSTGRES_SSLMODE" envDefault:"disable"`
	}
)

func NewConfig() (*Config, error) {
	cfg := &Config{}
	if err := env.Parse(cfg); err != nil {
		return nil, fmt.Errorf("%w: %w", ErrConfig, err)
	}

	return cfg, nil
}
