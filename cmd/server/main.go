package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/suapapa/croquis-king/internal/config"
)

func main() {
	cfg := config.Load()

	fmt.Printf("Starting server on port %d...\n", cfg.Port)

	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("OK"))
	})

	addr := fmt.Sprintf(":%d", cfg.Port)
	if err := http.ListenAndServe(addr, nil); err != nil {
		log.Fatalf("Server failed: %v", err)
	}
}
