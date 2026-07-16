package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/oni1997/onentry/services/api-go/config"
	"github.com/oni1997/onentry/services/api-go/crypto"
	"github.com/oni1997/onentry/services/api-go/database"
	"github.com/oni1997/onentry/services/api-go/handlers"
	"github.com/oni1997/onentry/services/api-go/middleware"
)

func main() {
	cfg := config.Load()
	db := database.MustNew(cfg.DatabaseURL)
	defer db.Close()

	cryptoClient := crypto.NewClient(cfg.CryptoServiceURL)

	rateLimiter := middleware.NewRateLimiter(cfg.RateLimitRPM)

	r := http.NewServeMux()

	r.Handle("/health", handlers.NewHealthHandler(db))
	r.Handle("/register", handlers.NewRegisterHandler(db, cryptoClient))
	r.Handle("/login", handlers.NewLoginHandler(db, cryptoClient, cfg))

	authMiddleware := middleware.AuthMiddleware(db)

	vaultHandler := handlers.NewVaultHandler(db, cryptoClient)
	r.Handle("/vault", authMiddleware(vaultHandler))

	generateHandler := handlers.NewGenerateHandler(db, cryptoClient)
	r.Handle("/generate", authMiddleware(generateHandler))

	meHandler := handlers.NewMeHandler(db)
	r.Handle("/me", authMiddleware(meHandler))

	handler := rateLimiter.Middleware(r)
	handler = corsMiddleware(handler)

	server := &http.Server{
		Addr:         ":" + cfg.Port,
		Handler:      loggingMiddleware(handler),
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	go func() {
		log.Printf("OnEntry API server starting on port %s", cfg.Port)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Server error: %v", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	_ = server.Shutdown(ctx)

	log.Println("Server stopped")
}

func loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		next.ServeHTTP(w, r)
		log.Printf("%s %s %s", r.Method, r.URL.Path, time.Since(start))
	})
}

func corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent)
			return
		}
		next.ServeHTTP(w, r)
	})
}
