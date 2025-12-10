package config

import "github.com/joho/godotenv"

type Config struct {
	HTTPConfig HTTPConfig
	PGConfig   PGConfig
}

func Load(path string) error {
	err := godotenv.Load(path)
	if err != nil {
		return err
	}

	return nil
}

func MustLoad(path string) *Config {

	err := Load(path)
	if err != nil {
		panic(err)
	}

	httpConfig, err := NewHTTPConfig()
	if err != nil {
		panic(err)
	}
	pgConfig, err := NewPGConfig()
	if err != nil {
		panic(err)
	}

	return &Config{
		HTTPConfig: httpConfig,
		PGConfig:   pgConfig,
	}
}
