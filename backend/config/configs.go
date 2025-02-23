package config

import (
	"os"
)

type (
	DB struct {
		DbHost     string
		DbUser     string
		DbName     string
		DbPassword string
	}
	HTTP struct {
		Port string
	}
	Redis struct {
		URL string
	}
	Configs struct {
		DB    *DB
		HTTP  *HTTP
		Redis *Redis
	}
)

func FromEnv() (*Configs, error) {
	return &Configs{
		&DB{
			os.Getenv("DB_HOST"),
			os.Getenv("DB_USER"),
			os.Getenv("DB_NAME"),
			os.Getenv("DB_PASSWORD"),
		},
		&HTTP{
			os.Getenv("HTTP_PORT"),
		},
		&Redis{
			os.Getenv("REDIS_URL"),
		},
	}, nil
}
