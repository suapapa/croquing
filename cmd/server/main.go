package main

import (
	"log"

	"github.com/suapapa/croquis-king/internal/config"
	"github.com/suapapa/croquis-king/internal/http"
	"github.com/suapapa/croquis-king/internal/lobby"
	"github.com/suapapa/croquis-king/internal/pixabay"
)

func main() {
	cfg := config.Load()

	store := lobby.NewMemoryStore()
	pixabayClient := pixabay.NewClient(cfg.PixabayAPIKey)
	srv, err := httpserver.New(cfg, store, pixabayClient)
	if err != nil {
		log.Fatalf("Server init failed: %v", err)
	}
	if err := srv.Run(); err != nil {
		log.Fatalf("Server failed: %v", err)
	}
}
