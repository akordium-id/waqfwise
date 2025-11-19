package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/akordium-id/waqfwise/internal/services/auth/handler"
	authRepo "github.com/akordium-id/waqfwise/internal/services/auth/repository"
	authService "github.com/akordium-id/waqfwise/internal/services/auth/service"
	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/rs/cors"
)

const version = "1.0.0"

func main() {
	log.Printf("🚀 WaqfWise Monolith v%s starting...", version)

	// Load configuration from environment
	config := loadConfig()

	// Connect to database
	db, err := connectDatabase(config.DatabaseURL)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	log.Println("✅ Connected to database")

	// Initialize services
	services := initializeServices(db, config)

	// Setup HTTP router
	router := setupRouter(services, config)

	// Create HTTP server
	srv := &http.Server{
		Addr:         ":" + config.Port,
		Handler:      router,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	// Start server in goroutine
	go func() {
		log.Printf("🌐 Server listening on port %s", config.Port)
		log.Printf("📋 Health check: http://localhost:%s/health", config.Port)
		log.Printf("📊 Metrics: http://localhost:%s/metrics", config.Port)
		log.Printf("🔐 Auth API: http://localhost:%s/api/v1/auth/*", config.Port)
		log.Printf("💰 Payment API: http://localhost:%s/api/v1/payments/*", config.Port)
		log.Printf("📢 Campaign API: http://localhost:%s/api/v1/campaigns/*", config.Port)
		log.Printf("🏛️  Asset API: http://localhost:%s/api/v1/assets/*", config.Port)

		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Failed to start server: %v", err)
		}
	}()

	// Graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("🛑 Shutting down server...")
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("Server forced to shutdown: %v", err)
	}

	log.Println("✅ Server exited gracefully")
}

// Config holds application configuration
type Config struct {
	Port        string
	DatabaseURL string
	JWTSecret   string
	Environment string
}

// loadConfig loads configuration from environment variables
func loadConfig() *Config {
	return &Config{
		Port:        getEnv("PORT", "8080"),
		DatabaseURL: getEnv("DATABASE_URL", "postgres://waqfwise:waqfwise@localhost:5432/waqfwise_community?sslmode=disable"),
		JWTSecret:   getEnv("JWT_SECRET", "your-secret-key-change-this-in-production"),
		Environment: getEnv("ENVIRONMENT", "development"),
	}
}

// Services holds all application services
type Services struct {
	AuthHandler *handler.Handler
	// PaymentHandler will be added when we implement it
	// CampaignHandler will be added when we implement it
	// AssetHandler will be added when we implement it
}

// initializeServices initializes all application services
func initializeServices(db *sql.DB, config *Config) *Services {
	// Initialize Auth service
	authRepository := authRepo.New(db)
	authSvc := authService.New(authRepository, config.JWTSecret)
	authHandler := handler.New(authSvc)

	// TODO: Initialize Payment service
	// paymentRepo := paymentRepo.New(db)
	// paymentSvc := paymentService.New(paymentRepo)
	// paymentHandler := paymentHandler.New(paymentSvc)

	// TODO: Initialize Campaign service
	// TODO: Initialize Asset service

	return &Services{
		AuthHandler: authHandler,
		// PaymentHandler: paymentHandler,
		// CampaignHandler: campaignHandler,
		// AssetHandler: assetHandler,
	}
}

// setupRouter sets up HTTP routes for all services
func setupRouter(services *Services, config *Config) http.Handler {
	router := mux.NewRouter()

	// Health check endpoint
	router.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, `{
			"status": "healthy",
			"service": "waqfwise-monolith",
			"version": "%s",
			"environment": "%s"
		}`, version, config.Environment)
	}).Methods("GET")

	// Metrics endpoint
	router.Handle("/metrics", promhttp.Handler())

	// API routes
	apiRouter := router.PathPrefix("/api/v1").Subrouter()

	// Auth routes
	authRouter := apiRouter.PathPrefix("/auth").Subrouter()
	services.AuthHandler.RegisterRoutes(authRouter)

	// Payment routes
	// paymentRouter := apiRouter.PathPrefix("/payments").Subrouter()
	// services.PaymentHandler.RegisterRoutes(paymentRouter)

	// Campaign routes
	// campaignRouter := apiRouter.PathPrefix("/campaigns").Subrouter()
	// services.CampaignHandler.RegisterRoutes(campaignRouter)

	// Asset routes
	// assetRouter := apiRouter.PathPrefix("/assets").Subrouter()
	// services.AssetHandler.RegisterRoutes(assetRouter)

	// CORS middleware
	c := cors.New(cors.Options{
		AllowedOrigins:   getAllowedOrigins(config.Environment),
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "PATCH", "OPTIONS"},
		AllowedHeaders:   []string{"*"},
		AllowCredentials: true,
		Debug:            config.Environment == "development",
	})

	// Logging middleware
	loggedRouter := loggingMiddleware(router)

	return c.Handler(loggedRouter)
}

// connectDatabase establishes database connection with connection pooling
func connectDatabase(url string) (*sql.DB, error) {
	db, err := sql.Open("postgres", url)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	// Configure connection pool
	db.SetMaxOpenConns(25)                 // Maximum open connections
	db.SetMaxIdleConns(5)                  // Maximum idle connections
	db.SetConnMaxLifetime(5 * time.Minute) // Connection lifetime

	// Test connection
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := db.PingContext(ctx); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	return db, nil
}

// loggingMiddleware logs HTTP requests
func loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		// Wrap response writer to capture status code
		wrapped := &responseWriter{ResponseWriter: w, statusCode: http.StatusOK}

		next.ServeHTTP(wrapped, r)

		// Skip logging for health and metrics
		if r.URL.Path == "/health" || r.URL.Path == "/metrics" {
			return
		}

		log.Printf(
			"%s %s %d %v %s",
			r.Method,
			r.URL.Path,
			wrapped.statusCode,
			time.Since(start),
			r.RemoteAddr,
		)
	})
}

// responseWriter wraps http.ResponseWriter to capture status code
type responseWriter struct {
	http.ResponseWriter
	statusCode int
}

func (rw *responseWriter) WriteHeader(code int) {
	rw.statusCode = code
	rw.ResponseWriter.WriteHeader(code)
}

// getAllowedOrigins returns allowed CORS origins based on environment
func getAllowedOrigins(env string) []string {
	if env == "production" {
		return []string{
			"https://waqfwise.id",
			"https://www.waqfwise.id",
			"https://app.waqfwise.id",
		}
	}
	// Development: allow all origins
	return []string{"*"}
}

// getEnv gets environment variable with fallback
func getEnv(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return fallback
}
