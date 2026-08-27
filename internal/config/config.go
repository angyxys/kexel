package config

import (
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	DATABASE_DSN string
	JWT_SECRET   string
	MAP_SECRET   string
}

func New() Config {
	godotenv.Load()
	rawDSN := os.Getenv("DATABASE_DSN")
	return Config{
		DATABASE_DSN: os.ExpandEnv(rawDSN),
		JWT_SECRET:   os.Getenv("JWT_SECRET"),
		MAP_SECRET:   os.Getenv("MAP_SECRET"),
	}
}
