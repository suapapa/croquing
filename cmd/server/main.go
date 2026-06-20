package main

import (
	"log"

	"github.com/suapapa/croquis-king/internal/config"
	"github.com/suapapa/croquis-king/internal/http"
)

func main() {
	cfg := config.Load()

	srv := httpserver.New(cfg)
	if err := srv.Run(); err != nil {
		log.Fatalf("Server failed: %v", err)
	}
}
