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
)

const shutdownTimeout = 10 * time.Second

// Server wraps the HTTP listener and Gin router.
type Server struct {
	cfg    *config.Config
	http   *http.Server
	router *gin.Engine
}

// New builds a Server from configuration.
func New(cfg *config.Config) *Server {
	router := newRouter()

	return &Server{
		cfg:    cfg,
		router: router,
		http: &http.Server{
			Addr:    fmt.Sprintf(":%d", cfg.Port),
			Handler: router,
		},
	}
}

// Run starts the HTTP server and blocks until it receives SIGINT or SIGTERM,
// then shuts down gracefully.
func (s *Server) Run() error {
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

	ctx, cancel := context.WithTimeout(context.Background(), shutdownTimeout)
	defer cancel()

	if err := s.http.Shutdown(ctx); err != nil {
		return fmt.Errorf("shutdown: %w", err)
	}

	log.Println("Server stopped")
	return nil
}
