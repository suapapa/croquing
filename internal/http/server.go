package httpserver

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"

	"github.com/suapapa/croquis-king/internal/config"
	"github.com/suapapa/croquis-king/internal/lobby"
	"github.com/suapapa/croquis-king/internal/pixabay"
	"github.com/suapapa/croquis-king/internal/timer"
	"github.com/suapapa/croquis-king/internal/ws"
)

const shutdownTimeout = 10 * time.Second

// Server wraps the HTTP listener and Gin router.
type Server struct {
	cfg       *config.Config
	http      *http.Server
	router    *gin.Engine
	lobbySync *ws.SnapshotSync
	scheduler *timer.Scheduler
}

// LobbySync returns the lobby snapshot broadcaster for state-changing handlers.
func (s *Server) LobbySync() *ws.SnapshotSync {
	return s.lobbySync
}

// New builds a Server from configuration and dependencies.
func New(cfg *config.Config, store lobby.Store, pixabayClient *pixabay.Client, lobbySync *ws.SnapshotSync, scheduler *timer.Scheduler) (*Server, error) {
	drawDuration, err := time.ParseDuration(cfg.DrawDuration)
	if err != nil {
		return nil, fmt.Errorf("parse DRAW_DURATION: %w", err)
	}

	var wsHandler *ws.Handler
	if lobbySync != nil {
		wsHandler = ws.NewHandler(lobbySync)
	}

	router := newRouter(store, drawDuration, pixabayClient, wsHandler, lobbySync)

	return &Server{
		cfg:       cfg,
		router:    router,
		lobbySync: lobbySync,
		scheduler: scheduler,
		http: &http.Server{
			Addr:    fmt.Sprintf(":%d", cfg.Port),
			Handler: router,
		},
	}, nil
}

// Run starts the HTTP server and blocks until it receives SIGINT or SIGTERM,
// then shuts down gracefully.
func (s *Server) Run() error {
	runCtx, stopScheduler := context.WithCancel(context.Background())
	defer stopScheduler()

	if s.scheduler != nil {
		go s.scheduler.Run(runCtx)
	}

	errCh := make(chan error, 1)

	go func() {
		log.Printf("HTTP server listening on %s", s.http.Addr)
		if err := s.http.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			errCh <- err
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	select {
	case err := <-errCh:
		return fmt.Errorf("listen: %w", err)
	case sig := <-quit:
		log.Printf("Received signal %s, shutting down...", sig)
	}

	stopScheduler()

	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), shutdownTimeout)
	defer shutdownCancel()

	if err := s.http.Shutdown(shutdownCtx); err != nil {
		return fmt.Errorf("shutdown: %w", err)
	}

	log.Println("Server stopped")
	return nil
}
