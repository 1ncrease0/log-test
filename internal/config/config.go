package config

import (
	"log"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	LogLevel  string `env:"LOG_LEVEL" env-default:"DEBUG"`
	DBAddress string `env:"DB_URL" env-required:"true"`
}

func MustLoad() Config {
	var cfg Config
	if err := cleanenv.ReadEnv(&cfg); err != nil {
		log.Fatalf("cannot read config: %s", err)
	}
	return cfg
}
