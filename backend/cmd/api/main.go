// Command api is the my-fitness-aibo HTTP backend.
package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/lionel/my-fitness-aibo/backend/internal/auth"
	"github.com/lionel/my-fitness-aibo/backend/internal/config"
	"github.com/lionel/my-fitness-aibo/backend/internal/db"
	"github.com/lionel/my-fitness-aibo/backend/internal/handlers"
	"github.com/lionel/my-fitness-aibo/backend/internal/store"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("config: %v", err)
	}

	ctx := context.Background()

	log.Println("running migrations...")
	if err := db.Migrate(cfg.DatabaseURL); err != nil {
		log.Fatalf("migrate: %v", err)
	}

	pool, err := db.Connect(ctx, cfg.DatabaseURL)
	if err != nil {
		log.Fatalf("connect db: %v", err)
	}
	defer pool.Close()

	tm := auth.NewTokenManager(cfg.JWTSecret, time.Duration(cfg.JWTTTLHours)*time.Hour)
	h := handlers.New(store.New(pool), tm)

	srv := &http.Server{
		Addr:              ":" + cfg.Port,
		Handler:           h.Router(cfg.CORSOrigins),
		ReadHeaderTimeout: 10 * time.Second,
	}

	go func() {
		log.Printf("listening on :%s", cfg.Port)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("server: %v", err)
		}
	}()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)
	<-stop

	log.Println("shutting down...")
	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := srv.Shutdown(shutdownCtx); err != nil {
		log.Printf("shutdown: %v", err)
	}
}
