package config

import (
	"log"

	"github.com/caarlosh/envconfig"
)

type Config struct {
	Port          int      `envconfig:"PORT" default:"8080"`
	PixabayAPIKey string   `envconfig:"PIXABAY_API_KEY" required:"true"`
	CORSOrigins   []string `envconfig:"CORS_ORIGINS" default:"*"`
	DrawDuration  string   `envconfig:"DRAW_DURATION" default:"5m"`
}

func Load() *Config {
	var cfg Config
	err := envconfig.Process("", &cfg)
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}
	return &cfg
}
