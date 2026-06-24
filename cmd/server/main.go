package main

import (
	"flag"
	"log/slog"
	"os"

	"github.com/suapapa/croquing/internal/config"
	httpserver "github.com/suapapa/croquing/internal/http"
	"github.com/suapapa/croquing/internal/logging"
	"github.com/suapapa/croquing/internal/lobby"
	"github.com/suapapa/croquing/internal/pixabay"
	"github.com/suapapa/croquing/internal/timer"
	"github.com/suapapa/croquing/internal/ws"
)

func main() {
	logFormat := flag.String("log-format", logging.FormatText, "log output format: text or json")
	flag.Parse()

	logging.Init(*logFormat)

	cfg, err := config.Load()
	if err != nil {
		slog.Error("failed to load config", "err", err)
		os.Exit(1)
	}

	store := lobby.NewMemoryStore()
	pixabayClient := pixabay.NewClient(cfg.PixabayAPIKey)
	lobbySync := ws.NewSnapshotSync(store)
	scheduler := timer.NewScheduler(store, lobbySync, timer.DefaultTickInterval)
	srv, err := httpserver.New(cfg, store, pixabayClient, lobbySync, scheduler)
	if err != nil {
		slog.Error("server init failed", "err", err)
		os.Exit(1)
	}

	slog.Info("starting server",
		"port", cfg.Port,
		"draw_duration", cfg.DrawDuration,
		"log_format", *logFormat,
	)

	if err := srv.Run(); err != nil {
		slog.Error("server failed", "err", err)
		os.Exit(1)
	}
}
