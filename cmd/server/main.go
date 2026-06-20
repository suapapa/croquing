package main

import (
	"log"

	"github.com/suapapa/croquis-king/internal/config"
	"github.com/suapapa/croquis-king/internal/http"
	"github.com/suapapa/croquis-king/internal/lobby"
	"github.com/suapapa/croquis-king/internal/pixabay"
	"github.com/suapapa/croquis-king/internal/timer"
	"github.com/suapapa/croquis-king/internal/ws"
)

func main() {
	cfg := config.Load()

	store := lobby.NewMemoryStore()
	pixabayClient := pixabay.NewClient(cfg.PixabayAPIKey)
	lobbySync := ws.NewSnapshotSync(store)
	scheduler := timer.NewScheduler(store, lobbySync, timer.DefaultTickInterval)
	srv, err := httpserver.New(cfg, store, pixabayClient, lobbySync, scheduler)
	if err != nil {
		log.Fatalf("Server init failed: %v", err)
	}
	if err := srv.Run(); err != nil {
		log.Fatalf("Server failed: %v", err)
	}
}
