package config

import (
	"github.com/caarlos0/env/v6"
	"github.com/joho/godotenv"
	"github.com/mcuadros/go-defaults"
)

type Config struct {
	HTTPServer HTTPServerConfig
	DB         DBConfig
}

type HTTPServerConfig struct {
	Host string `env:"SERVER_HOST" default:"localhost"`
	Port int    `env:"SERVER_PORT" default:"8080"`
}

type DBConfig struct {
	DSN string `env:"DATABASE_DSN"`
}

func NewConfig(filenames ...string) (*Config, error) {
	_ = godotenv.Load(filenames...)
	cfg := &Config{}
	if err := env.Parse(cfg); err != nil {
		return nil, err
	}
	defaults.SetDefaults(cfg)
	return cfg, nil
}
