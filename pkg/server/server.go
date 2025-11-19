package server

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/akordium-id/waqfwise/pkg/auth"
	"github.com/akordium-id/waqfwise/pkg/cache"
	"github.com/akordium-id/waqfwise/pkg/config"
	"github.com/akordium-id/waqfwise/pkg/database"
	"github.com/akordium-id/waqfwise/pkg/logger"
	"github.com/akordium-id/waqfwise/pkg/metrics"
	"github.com/akordium-id/waqfwise/pkg/middleware"
	"github.com/akordium-id/waqfwise/pkg/payment"
	"github.com/akordium-id/waqfwise/pkg/queue"
	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"go.uber.org/zap"
)

// Server represents the HTTP server
type Server struct {
	config         *config.Config
	router         *gin.Engine
	httpServer     *http.Server
	db             *database.DB
	redis          *cache.Redis
	kafka          *queue.Kafka
	logger         *zap.Logger
	metrics        *metrics.Metrics
	jwtManager     *auth.JWTManager
	oauth2Manager  *auth.OAuth2Manager
	sessionStore   *cache.SessionStore
	rateLimiter    *cache.RateLimiter
	midtransGateway *payment.MidtransGateway
	xenditGateway   *payment.XenditGateway
}

// New creates a new server instance
func New(cfg *config.Config) (*Server, error) {
	// Initialize logger
	log, err := logger.New(&cfg.Logging)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize logger: %w", err)
	}

	// Initialize database
	db, err := database.New(&cfg.Database, log)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize database: %w", err)
	}

	// Initialize Redis
	redis, err := cache.New(&cfg.Redis, log)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize Redis: %w", err)
	}

	// Initialize Kafka
	kafka, err := queue.New(&cfg.Kafka, log)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize Kafka: %w", err)
	}

	// Initialize metrics
	metricsInstance := metrics.New("waqfwise")

	// Initialize JWT manager
	jwtManager := auth.NewJWTManager(&cfg.JWT)

	// Initialize OAuth2 manager
	oauth2Manager := auth.NewOAuth2Manager(&cfg.OAuth2)

	// Initialize session store
	sessionStore := cache.NewSessionStore(redis, time.Duration(cfg.Session.MaxAge)*time.Second)

	// Initialize rate limiter
	rateLimiter := cache.NewRateLimiter(redis)

	// Initialize payment gateways
	midtransGateway := payment.NewMidtransGateway(&cfg.Midtrans, log)
	xenditGateway := payment.NewXenditGateway(&cfg.Xendit, log)

	// Set Gin mode
	if cfg.App.IsProduction() {
		gin.SetMode(gin.ReleaseMode)
	}

	// Create router
	router := gin.New()

	// Create server
	server := &Server{
		config:          cfg,
		router:          router,
		db:              db,
		redis:           redis,
		kafka:           kafka,
		logger:          log,
		metrics:         metricsInstance,
		jwtManager:      jwtManager,
		oauth2Manager:   oauth2Manager,
		sessionStore:    sessionStore,
		rateLimiter:     rateLimiter,
		midtransGateway: midtransGateway,
		xenditGateway:   xenditGateway,
	}

	// Setup middleware
	server.setupMiddleware()

	// Setup routes
	server.setupRoutes()

	return server, nil
}

// setupMiddleware sets up middleware
func (s *Server) setupMiddleware() {
	// Recovery middleware
	s.router.Use(gin.Recovery())

	// Logger middleware
	s.router.Use(middleware.Logger(s.logger))

	// Metrics middleware
	s.router.Use(middleware.Metrics(s.metrics))

	// CORS middleware
	s.router.Use(middleware.CORS(s.config.CORS))

	// Rate limiting middleware (if enabled)
	if s.config.RateLimit.Enabled {
		s.router.Use(middleware.RateLimit(s.rateLimiter, s.config.RateLimit.RequestsPerMinute))
	}
}

// setupRoutes sets up routes
func (s *Server) setupRoutes() {
	// Health check
	s.router.GET("/health", s.healthHandler)
	s.router.GET("/ready", s.readyHandler)

	// Metrics endpoint
	if s.config.Metrics.Enabled {
		s.router.GET("/metrics", gin.WrapH(promhttp.Handler()))
	}

	// API v1
	v1 := s.router.Group("/api/v1")
	{
		// Public routes
		public := v1.Group("/public")
		{
			public.GET("/ping", s.pingHandler)
		}

		// Auth routes
		authRoutes := v1.Group("/auth")
		{
			authRoutes.POST("/register", s.registerHandler)
			authRoutes.POST("/login", s.loginHandler)
			authRoutes.POST("/refresh", s.refreshTokenHandler)

			// OAuth2 routes
			authRoutes.GET("/google", s.googleAuthHandler)
			authRoutes.GET("/google/callback", s.googleCallbackHandler)
			authRoutes.GET("/facebook", s.facebookAuthHandler)
			authRoutes.GET("/facebook/callback", s.facebookCallbackHandler)
		}

		// Protected routes
		protected := v1.Group("")
		protected.Use(auth.JWTMiddleware(s.jwtManager))
		{
			protected.GET("/profile", s.profileHandler)
			protected.POST("/logout", s.logoutHandler)
		}
	}
}

// Start starts the HTTP server
func (s *Server) Start() error {
	addr := fmt.Sprintf(":%d", s.config.App.HTTPPort)

	s.httpServer = &http.Server{
		Addr:           addr,
		Handler:        s.router,
		ReadTimeout:    15 * time.Second,
		WriteTimeout:   15 * time.Second,
		IdleTimeout:    60 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}

	s.logger.Info("starting HTTP server",
		zap.String("address", addr),
		zap.String("environment", s.config.App.Environment),
	)

	if err := s.httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		return fmt.Errorf("failed to start HTTP server: %w", err)
	}

	return nil
}

// Shutdown gracefully shuts down the server
func (s *Server) Shutdown(ctx context.Context) error {
	s.logger.Info("shutting down server...")

	// Shutdown HTTP server
	if err := s.httpServer.Shutdown(ctx); err != nil {
		s.logger.Error("HTTP server shutdown error", zap.Error(err))
	}

	// Close Kafka
	if err := s.kafka.Close(); err != nil {
		s.logger.Error("Kafka shutdown error", zap.Error(err))
	}

	// Close Redis
	if err := s.redis.Close(); err != nil {
		s.logger.Error("Redis shutdown error", zap.Error(err))
	}

	// Close database
	if err := s.db.Close(); err != nil {
		s.logger.Error("database shutdown error", zap.Error(err))
	}

	s.logger.Info("server shutdown complete")
	return nil
}

// Handler placeholders - these will be implemented in actual handlers

func (s *Server) healthHandler(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"status": "healthy"})
}

func (s *Server) readyHandler(c *gin.Context) {
	// Check database
	if err := s.db.HealthCheck(); err != nil {
		c.JSON(http.StatusServiceUnavailable, gin.H{"ready": false, "error": "database not ready"})
		return
	}

	// Check Redis
	if err := s.redis.HealthCheck(); err != nil {
		c.JSON(http.StatusServiceUnavailable, gin.H{"ready": false, "error": "redis not ready"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"ready": true})
}

func (s *Server) pingHandler(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "pong"})
}

// Placeholder handlers
func (s *Server) registerHandler(c *gin.Context)          { c.JSON(http.StatusNotImplemented, gin.H{"error": "not implemented"}) }
func (s *Server) loginHandler(c *gin.Context)             { c.JSON(http.StatusNotImplemented, gin.H{"error": "not implemented"}) }
func (s *Server) refreshTokenHandler(c *gin.Context)      { c.JSON(http.StatusNotImplemented, gin.H{"error": "not implemented"}) }
func (s *Server) googleAuthHandler(c *gin.Context)        { c.JSON(http.StatusNotImplemented, gin.H{"error": "not implemented"}) }
func (s *Server) googleCallbackHandler(c *gin.Context)    { c.JSON(http.StatusNotImplemented, gin.H{"error": "not implemented"}) }
func (s *Server) facebookAuthHandler(c *gin.Context)      { c.JSON(http.StatusNotImplemented, gin.H{"error": "not implemented"}) }
func (s *Server) facebookCallbackHandler(c *gin.Context)  { c.JSON(http.StatusNotImplemented, gin.H{"error": "not implemented"}) }
func (s *Server) profileHandler(c *gin.Context)           { c.JSON(http.StatusNotImplemented, gin.H{"error": "not implemented"}) }
func (s *Server) logoutHandler(c *gin.Context)            { c.JSON(http.StatusNotImplemented, gin.H{"error": "not implemented"}) }
