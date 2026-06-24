package config

import (
	"fmt"

	"github.com/kelseyhightower/envconfig"
)

type Config struct {
	Port          int      `envconfig:"PORT" default:"8080"`
	PixabayAPIKey string   `envconfig:"PIXABAY_API_KEY" required:"true"`
	CORSOrigins   []string `envconfig:"CORS_ORIGINS" default:"*"`
	DrawDuration  string   `envconfig:"DRAW_DURATION" default:"5m"`
	AppName       string   `envconfig:"APP_NAME" default:"Croquing"`
}

func Load() (*Config, error) {
	var cfg Config
	if err := envconfig.Process("", &cfg); err != nil {
		return nil, fmt.Errorf("load config: %w", err)
	}
	return &cfg, nil
}
