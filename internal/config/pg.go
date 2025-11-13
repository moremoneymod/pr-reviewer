package config

import (
	"errors"
	"os"
)

const (
	pgDSNEnvName = "PG_DSN"
)

type PGConfig struct {
	dsn string
}

func NewPGConfig() (*PGConfig, error) {
	dsn := os.Getenv(pgDSNEnvName)
	if len(dsn) == 0 {
		return nil, errors.New("pg dsn not found")
	}

	return &PGConfig{
		dsn: dsn,
	}, nil
}

func (cfg *PGConfig) DSN() string {
	return cfg.dsn
}
