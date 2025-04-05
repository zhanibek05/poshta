package config

import (
	"github.com/caarlos0/env/v6"
	"github.com/joho/godotenv"
	"github.com/mcuadros/go-defaults"

	"time"
)

type Config struct {
	HTTPServer HTTPServerConfig
	DB         DBConfig
	JWT 	   JWTConfig
}

type HTTPServerConfig struct {
	Host string `env:"SERVER_HOST" default:"localhost"`
	Port int    `env:"SERVER_PORT" default:"8080"`
}

type DBConfig struct {
	DSN string `env:"DATABASE_DSN"`
}

type JWTConfig struct {
	SecretKey       string        `env:"JWT_SECRET_KEY" default:"your_secret_key_here"`
	AccessTokenTTL  time.Duration `env:"JWT_ACCESS_TOKEN_TTL" default:"15m"`
	RefreshTokenTTL time.Duration `env:"JWT_REFRESH_TOKEN_TTL" default:"72h"`
	Issuer          string        `env:"JWT_ISSUER" default:"poshta-app"`
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
