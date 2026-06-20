package main

import (
	"log"

	"github.com/suapapa/croquis-king/internal/config"
	"github.com/suapapa/croquis-king/internal/http"
	"github.com/suapapa/croquis-king/internal/lobby"
)

func main() {
	cfg := config.Load()

	store := lobby.NewMemoryStore()
	srv, err := httpserver.New(cfg, store)
	if err != nil {
		log.Fatalf("Server init failed: %v", err)
	}
	if err := srv.Run(); err != nil {
		log.Fatalf("Server failed: %v", err)
	}
}
