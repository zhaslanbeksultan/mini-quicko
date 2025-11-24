package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"mini-quicko/internal/config"
	"mini-quicko/internal/handlers"
	"mini-quicko/internal/service"
	"mini-quicko/internal/storage"
)

func main() {
	cfg := config.Load()

	store := storage.NewMemoryStorage()
	kaspiSvc := service.NewKaspiService(cfg.MockDataPath)
	analyzerSvc := service.NewAnalyzer(store, kaspiSvc)

	handler := handlers.NewHandler(analyzerSvc)
	router := handlers.SetupRoutes(handler)

	srv := &http.Server{
		Addr:         cfg.ServerAddress,
		Handler:      router,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	go func() {
		log.Printf("Starting server on %s", cfg.ServerAddress)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Server error: %v", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("Shutting down server...")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("Server forced to shutdown: %v", err)
	}

	log.Println("Server stopped")
}
