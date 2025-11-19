package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/gorilla/mux"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/rs/cors"
)

type ServiceConfig struct {
	Name string
	URL  string
}

var services = map[string]ServiceConfig{
	"auth":        {Name: "auth-service", URL: getEnv("AUTH_SERVICE_URL", "http://localhost:8001")},
	"payment":     {Name: "payment-service", URL: getEnv("PAYMENT_SERVICE_URL", "http://localhost:8002")},
	"campaign":    {Name: "campaign-service", URL: getEnv("CAMPAIGN_SERVICE_URL", "http://localhost:8003")},
	"asset":       {Name: "asset-service", URL: getEnv("ASSET_SERVICE_URL", "http://localhost:8004")},
	"analytics":   {Name: "analytics-service", URL: getEnv("ANALYTICS_SERVICE_URL", "http://localhost:8005")},
	"integration": {Name: "integration-service", URL: getEnv("INTEGRATION_SERVICE_URL", "http://localhost:8006")},
}

func main() {
	port := getEnv("PORT", "8000")

	router := mux.NewRouter()

	// Health check
	router.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, `{"status":"healthy","service":"api-gateway"}`)
	}).Methods("GET")

	// Metrics
	router.Handle("/metrics", promhttp.Handler())

	// Service routes
	apiRouter := router.PathPrefix("/api/v1").Subrouter()

	// Auth service routes
	apiRouter.PathPrefix("/auth").HandlerFunc(createProxyHandler("auth"))

	// Payment service routes
	apiRouter.PathPrefix("/payments").HandlerFunc(createProxyHandler("payment"))
	apiRouter.PathPrefix("/donations").HandlerFunc(createProxyHandler("payment"))

	// Campaign service routes
	apiRouter.PathPrefix("/campaigns").HandlerFunc(createProxyHandler("campaign"))

	// Asset service routes
	apiRouter.PathPrefix("/assets").HandlerFunc(createProxyHandler("asset"))

	// Analytics service routes (Enterprise only)
	apiRouter.PathPrefix("/analytics").HandlerFunc(createProxyHandler("analytics"))

	// Integration service routes (Enterprise only)
	apiRouter.PathPrefix("/integrations").HandlerFunc(createProxyHandler("integration"))
	apiRouter.PathPrefix("/webhooks").HandlerFunc(createProxyHandler("integration"))

	// CORS middleware
	c := cors.New(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "PATCH", "OPTIONS"},
		AllowedHeaders:   []string{"*"},
		AllowCredentials: true,
		Debug:            false,
	})

	// Logging middleware
	loggedRouter := loggingMiddleware(router)

	// Create HTTP server
	srv := &http.Server{
		Addr:         ":" + port,
		Handler:      c.Handler(loggedRouter),
		ReadTimeout:  30 * time.Second,
		WriteTimeout: 30 * time.Second,
		IdleTimeout:  120 * time.Second,
	}

	// Start server
	go func() {
		log.Printf("API Gateway starting on port %s", port)
		log.Println("Service routing:")
		for path, svc := range services {
			log.Printf("  /%s -> %s", path, svc.URL)
		}

		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Failed to start server: %v", err)
		}
	}()

	// Graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("Shutting down API Gateway...")
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("API Gateway forced to shutdown: %v", err)
	}

	log.Println("API Gateway exited")
}

// createProxyHandler creates a reverse proxy handler for a service
func createProxyHandler(serviceName string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		service, ok := services[serviceName]
		if !ok {
			http.Error(w, "Service not found", http.StatusNotFound)
			return
		}

		targetURL, err := url.Parse(service.URL)
		if err != nil {
			log.Printf("Failed to parse service URL: %v", err)
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}

		// Create reverse proxy
		proxy := httputil.NewSingleHostReverseProxy(targetURL)

		// Modify request
		originalPath := r.URL.Path
		r.URL.Host = targetURL.Host
		r.URL.Scheme = targetURL.Scheme
		r.Header.Set("X-Forwarded-Host", r.Header.Get("Host"))
		r.Header.Set("X-Origin-Host", targetURL.Host)
		r.Host = targetURL.Host

		// Custom error handler
		proxy.ErrorHandler = func(w http.ResponseWriter, r *http.Request, err error) {
			log.Printf("Proxy error for %s: %v", service.Name, err)
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusBadGateway)
			fmt.Fprintf(w, `{"error":"Service unavailable","service":"%s"}`, service.Name)
		}

		log.Printf("Proxying: %s %s -> %s%s", r.Method, originalPath, service.URL, r.URL.Path)
		proxy.ServeHTTP(w, r)
	}
}

// loggingMiddleware logs all HTTP requests
func loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		// Skip logging for health and metrics endpoints
		if !strings.Contains(r.URL.Path, "/health") && !strings.Contains(r.URL.Path, "/metrics") {
			defer func() {
				log.Printf(
					"%s %s %s %v",
					r.Method,
					r.RequestURI,
					r.RemoteAddr,
					time.Since(start),
				)
			}()
		}

		next.ServeHTTP(w, r)
	})
}

func getEnv(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return fallback
}
