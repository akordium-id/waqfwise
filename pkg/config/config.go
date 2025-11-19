package config

import (
	"fmt"
	"time"

	"github.com/spf13/viper"
)

// Config holds all configuration for the application
type Config struct {
	App        AppConfig
	Database   DatabaseConfig
	Redis      RedisConfig
	Kafka      KafkaConfig
	JWT        JWTConfig
	OAuth2     OAuth2Config
	Midtrans   MidtransConfig
	Xendit     XenditConfig
	Logging    LoggingConfig
	Metrics    MetricsConfig
	RateLimit  RateLimitConfig
	CORS       CORSConfig
	Session    SessionConfig
	Enterprise EnterpriseConfig
}

// AppConfig holds application configuration
type AppConfig struct {
	Name        string
	Version     string
	Environment string
	HTTPPort    int
	HTTPSPort   int
}

// DatabaseConfig holds database configuration
type DatabaseConfig struct {
	Host            string
	Port            int
	User            string
	Password        string
	Name            string
	SSLMode         string
	MaxOpenConns    int
	MaxIdleConns    int
	ConnMaxLifetime time.Duration
}

// RedisConfig holds Redis configuration
type RedisConfig struct {
	Host       string
	Port       int
	Password   string
	DB         int
	MaxRetries int
	PoolSize   int
}

// KafkaConfig holds Kafka configuration
type KafkaConfig struct {
	Brokers           []string
	GroupID           string
	AutoOffsetReset   string
	EnableAutoCommit  bool
	SessionTimeout    time.Duration
	HeartbeatInterval time.Duration
}

// JWTConfig holds JWT configuration
type JWTConfig struct {
	Secret              string
	AccessTokenExpire   time.Duration
	RefreshTokenExpire  time.Duration
	Issuer              string
}

// OAuth2Config holds OAuth2 configuration
type OAuth2Config struct {
	Google   OAuth2ProviderConfig
	Facebook OAuth2ProviderConfig
}

// OAuth2ProviderConfig holds OAuth2 provider configuration
type OAuth2ProviderConfig struct {
	ClientID     string
	ClientSecret string
	RedirectURL  string
}

// MidtransConfig holds Midtrans payment gateway configuration
type MidtransConfig struct {
	ServerKey   string
	ClientKey   string
	Environment string // sandbox or production
	MerchantID  string
}

// XenditConfig holds Xendit payment gateway configuration
type XenditConfig struct {
	SecretKey    string
	WebhookToken string
	Environment  string // sandbox or production
}

// LoggingConfig holds logging configuration
type LoggingConfig struct {
	Level  string
	Format string // json or console
	Output string // stdout, stderr, or file path
}

// MetricsConfig holds metrics configuration
type MetricsConfig struct {
	Enabled bool
	Port    int
}

// RateLimitConfig holds rate limiting configuration
type RateLimitConfig struct {
	Enabled           bool
	RequestsPerMinute int
}

// CORSConfig holds CORS configuration
type CORSConfig struct {
	AllowedOrigins []string
	AllowedMethods []string
	AllowedHeaders []string
}

// SessionConfig holds session configuration
type SessionConfig struct {
	Secret string
	MaxAge int // in seconds
}

// EnterpriseConfig holds enterprise-specific configuration
type EnterpriseConfig struct {
	LicenseKey          string
	MultiTenancyEnabled bool
}

// Load loads configuration from environment variables and config files
func Load(configPath string) (*Config, error) {
	v := viper.New()

	// Set default values
	setDefaults(v)

	// Read from environment variables
	v.AutomaticEnv()

	// Read from config file if provided
	if configPath != "" {
		v.SetConfigFile(configPath)
		if err := v.ReadInConfig(); err != nil {
			return nil, fmt.Errorf("failed to read config file: %w", err)
		}
	}

	var config Config
	if err := v.Unmarshal(&config); err != nil {
		return nil, fmt.Errorf("failed to unmarshal config: %w", err)
	}

	return &config, nil
}

// setDefaults sets default configuration values
func setDefaults(v *viper.Viper) {
	// App defaults
	v.SetDefault("app.name", "WaqfWise")
	v.SetDefault("app.version", "1.0.0")
	v.SetDefault("app.environment", "development")
	v.SetDefault("app.httpport", 8080)
	v.SetDefault("app.httpsport", 8443)

	// Database defaults
	v.SetDefault("database.host", "localhost")
	v.SetDefault("database.port", 5432)
	v.SetDefault("database.user", "waqfwise")
	v.SetDefault("database.password", "waqfwise_dev_password")
	v.SetDefault("database.name", "waqfwise")
	v.SetDefault("database.sslmode", "disable")
	v.SetDefault("database.maxopenconns", 25)
	v.SetDefault("database.maxidleconns", 5)
	v.SetDefault("database.connmaxlifetime", "5m")

	// Redis defaults
	v.SetDefault("redis.host", "localhost")
	v.SetDefault("redis.port", 6379)
	v.SetDefault("redis.password", "")
	v.SetDefault("redis.db", 0)
	v.SetDefault("redis.maxretries", 3)
	v.SetDefault("redis.poolsize", 10)

	// Kafka defaults
	v.SetDefault("kafka.brokers", []string{"localhost:29092"})
	v.SetDefault("kafka.groupid", "waqfwise-consumer-group")
	v.SetDefault("kafka.autooffsetreset", "earliest")
	v.SetDefault("kafka.enableautocommit", true)
	v.SetDefault("kafka.sessiontimeout", "10s")
	v.SetDefault("kafka.heartbeatinterval", "3s")

	// JWT defaults
	v.SetDefault("jwt.secret", "change-this-secret-in-production")
	v.SetDefault("jwt.accesstokenexpire", "15m")
	v.SetDefault("jwt.refreshtokenexpire", "168h")
	v.SetDefault("jwt.issuer", "waqfwise")

	// Logging defaults
	v.SetDefault("logging.level", "info")
	v.SetDefault("logging.format", "json")
	v.SetDefault("logging.output", "stdout")

	// Metrics defaults
	v.SetDefault("metrics.enabled", true)
	v.SetDefault("metrics.port", 9091)

	// Rate limiting defaults
	v.SetDefault("ratelimit.enabled", true)
	v.SetDefault("ratelimit.requestsperminute", 60)

	// CORS defaults
	v.SetDefault("cors.allowedorigins", []string{"http://localhost:3000"})
	v.SetDefault("cors.allowedmethods", []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"})
	v.SetDefault("cors.allowedheaders", []string{"Origin", "Content-Type", "Accept", "Authorization"})

	// Session defaults
	v.SetDefault("session.secret", "change-this-secret-in-production")
	v.SetDefault("session.maxage", 86400)

	// Enterprise defaults
	v.SetDefault("enterprise.licensekey", "")
	v.SetDefault("enterprise.multitenancyenabled", false)
}

// GetDSN returns the database connection string
func (c *DatabaseConfig) GetDSN() string {
	return fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
		c.Host, c.Port, c.User, c.Password, c.Name, c.SSLMode)
}

// GetRedisAddr returns the Redis address
func (c *RedisConfig) GetRedisAddr() string {
	return fmt.Sprintf("%s:%d", c.Host, c.Port)
}

// IsSandbox returns true if Midtrans is in sandbox mode
func (c *MidtransConfig) IsSandbox() bool {
	return c.Environment == "sandbox"
}

// IsSandbox returns true if Xendit is in sandbox mode
func (c *XenditConfig) IsSandbox() bool {
	return c.Environment == "sandbox"
}

// IsProduction returns true if the app is in production mode
func (c *AppConfig) IsProduction() bool {
	return c.Environment == "production"
}
