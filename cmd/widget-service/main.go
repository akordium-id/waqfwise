package main

import (
	"context"
	"embed"
	"io/fs"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/akordium-id/waqfwise/internal/shared/config"
	"github.com/akordium-id/waqfwise/internal/widget/handler"
	"github.com/akordium-id/waqfwise/internal/widget/repository"
	"github.com/akordium-id/waqfwise/internal/widget/service"
	"github.com/gin-gonic/gin"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

//go:embed web/widget
var widgetFS embed.FS

func main() {
	cfg := config.LoadWidgetConfig()

	// Connect to database
	db, err := gorm.Open(postgres.Open(cfg.DatabaseURL), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Warn),
	})
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	// Wire dependencies
	repo := repository.New(db)
	svc := service.New(repo, cfg.BaseURL)
	h := handler.New(svc)

	// Setup router
	if cfg.LogLevel != "debug" {
		gin.SetMode(gin.ReleaseMode)
	}
	r := gin.New()
	r.Use(gin.Recovery())

	// Serve static assets
	subFS, err := fs.Sub(widgetFS, "web/widget")
	if err != nil {
		log.Fatalf("Failed to create sub filesystem: %v", err)
	}
	r.StaticFS("/widget", http.FS(subFS))

	// Widget page
	r.GET("/widget", h.WidgetPage)

	// Public API
	r.GET("/api/v1/widget/campaigns", h.GetCampaigns)
	r.GET("/api/v1/mosques/:slug/config", h.GetMosqueConfig)

	// Health check
	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "healthy"})
	})

	// Start server
	srv := &http.Server{
		Addr:         ":" + cfg.Port,
		Handler:      r,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	go func() {
		log.Printf("Widget service starting on port %s", cfg.Port)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Server failed: %v", err)
		}
	}()

	// Graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Shutting down...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("Server forced to shutdown: %v", err)
	}
	log.Println("Shutdown complete")
}
