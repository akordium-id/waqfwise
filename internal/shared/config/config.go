package config

import "os"

// WidgetConfig holds widget service configuration
type WidgetConfig struct {
	Port        string
	DatabaseURL string
	LogLevel    string
	BaseURL     string
}

// LoadWidgetConfig loads widget configuration from environment
func LoadWidgetConfig() *WidgetConfig {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	baseURL := os.Getenv("WIDGET_BASE_URL")
	if baseURL == "" {
		baseURL = "http://localhost:8080"
	}
	logLevel := os.Getenv("LOG_LEVEL")
	if logLevel == "" {
		logLevel = "info"
	}
	return &WidgetConfig{
		Port:        port,
		DatabaseURL: os.Getenv("DATABASE_URL"),
		LogLevel:    logLevel,
		BaseURL:     baseURL,
	}
}
