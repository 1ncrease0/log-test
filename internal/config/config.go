package config

import (
	"log"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	LogLevel   string
	DBAddress  string
	HTTPConfig HTTPConfig
}

type HTTPConfig struct {
	Port string
}

type envConfig struct {
	LogLevel    string `env:"LOG_LEVEL" env-default:"INFO"`
	DatabaseURL string `env:"DATABASE_URL" env-required:"true"`
	Port        string `env:"PORT" env-default:"8080"`
}

func MustLoad() Config {
	var raw envConfig
	if err := cleanenv.ReadEnv(&raw); err != nil {
		log.Fatalf("cannot read config: %s", err)
	}
	return Config{
		LogLevel:  raw.LogLevel,
		DBAddress: raw.DatabaseURL,
		HTTPConfig: HTTPConfig{
			Port: raw.Port,
		},
	}
}
