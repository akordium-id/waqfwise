package handler

import (
	"encoding/json"
	"errors"
	"html"
	"net/http"
	"net/url"
	"regexp"
	"strconv"
	"strings"

	"github.com/akordium-id/waqfwise/internal/widget/repository"
	"github.com/akordium-id/waqfwise/internal/widget/service"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// Handler serves widget HTTP endpoints
type Handler struct {
	svc *service.Service
}

// New creates a widget handler
func New(svc *service.Service) *Handler {
	return &Handler{svc: svc}
}

// WidgetPage renders the iframe HTML
func (h *Handler) WidgetPage(c *gin.Context) {
	mosqueID := c.Query("mosque_id")
	if mosqueID == "" {
		c.String(http.StatusBadRequest, "mosque_id required")
		return
	}

	mosque, _, err := h.svc.GetMosqueConfig(c.Request.Context(), mosqueID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.String(http.StatusNotFound, "Mosque not found")
			return
		}
		c.String(http.StatusInternalServerError, "Internal error")
		return
	}

	// Referrer validation (host comparison, not prefix)
	if mosque.Domain != "" {
		referer := c.Request.Referer()
		if referer != "" {
			refURL, err := url.Parse(referer)
			if err != nil {
				c.String(http.StatusForbidden, "Invalid referrer")
				return
			}
			domainURL, err := url.Parse(mosque.Domain)
			if err != nil {
				c.String(http.StatusForbidden, "Invalid referrer")
				return
			}
			if refURL.Host != domainURL.Host {
				c.String(http.StatusForbidden, "Invalid referrer")
				return
			}
		}
	}

	// Render inline HTML with config
	primaryColor := c.Query("primary_color")
	if primaryColor == "" {
		primaryColor = strings.TrimPrefix(mosque.PrimaryColor, "#")
	}
	if !validHexColor(primaryColor) {
		c.String(http.StatusBadRequest, "Invalid primary_color")
		return
	}

	logoURL := c.Query("logo_url")
	if logoURL == "" {
		logoURL = mosque.LogoURL
	}

	theme := c.Query("theme")
	if theme == "" {
		theme = "light"
	}
	if theme != "light" && theme != "dark" {
		c.String(http.StatusBadRequest, "Invalid theme")
		return
	}

	limit := c.Query("limit")
	if limit == "" {
		limit = "6"
	}
	if _, err := strconv.Atoi(limit); err != nil {
		c.String(http.StatusBadRequest, "Invalid limit")
		return
	}

	types := c.Query("types")

	html := renderWidgetHTML(mosque.Name, primaryColor, logoURL, theme, mosqueID, limit, types)
	c.Data(http.StatusOK, "text/html; charset=utf-8", []byte(html))
}

// validHexColor reports whether s is a 6-digit hex color (with or without leading #).
func validHexColor(s string) bool {
	s = strings.TrimPrefix(s, "#")
	matched, _ := regexp.MatchString(`^[0-9a-fA-F]{6}$`, s)
	return matched
}

// renderWidgetHTML builds the widget HTML shell
func renderWidgetHTML(mosqueName, primaryColor, logoURL, theme, mosqueID, limit, types string) string {
	safeName := html.EscapeString(mosqueName)
	safeLogo := html.EscapeString(logoURL)

	type widgetConfig struct {
		MosqueID     string `json:"mosqueId"`
		PrimaryColor string `json:"primaryColor"`
		LogoURL      string `json:"logoUrl"`
		Limit        int    `json:"limit"`
		Types        string `json:"types"`
		Theme        string `json:"theme"`
		APIBase      string `json:"apiBase"`
	}

	limitInt, _ := strconv.Atoi(limit)
	cfg := widgetConfig{
		MosqueID:     mosqueID,
		PrimaryColor: "#" + primaryColor,
		LogoURL:      logoURL,
		Limit:        limitInt,
		Types:        types,
		Theme:        theme,
		APIBase:      "/api/v1/widget",
	}
	jsConfig, _ := json.Marshal(cfg)

	return `<!DOCTYPE html>
<html lang="id" data-theme="` + html.EscapeString(theme) + `">
<head>
<meta charset="UTF-8">
<meta name="viewport" content="width=device-width, initial-scale=1.0">
<title>` + safeName + ` - Kampanye Wakaf</title>
<style>:root{--waqf-primary:#` + html.EscapeString(primaryColor) + `;--waqf-logo:url(` + safeLogo + `)}</style>
<link rel="stylesheet" href="/widget/styles.css">
</head>
<body>
<div id="widget-root">
  <div class="waqf-skeleton" id="skeleton"><div class="waqf-skeleton-card"></div><div class="waqf-skeleton-card"></div><div class="waqf-skeleton-card"></div></div>
  <div class="waqf-error hidden" id="error"><p>Gagal memuat kampanye</p><button onclick="loadCampaigns()">Coba Lagi</button></div>
  <div class="waqf-empty hidden" id="empty"><p>Belum ada kampanye aktif</p></div>
  <div class="waqf-grid hidden" id="grid"></div>
</div>
<script>
window.WIDGET_CONFIG = ` + string(jsConfig) + `;
</script>
<script src="/widget/app.js"></script>
</body>
</html>`
}

// GetCampaigns returns public campaign data
func (h *Handler) GetCampaigns(c *gin.Context) {
	mosqueSlug := c.Query("mosque_slug")
	if mosqueSlug == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "mosque_slug required"})
		return
	}

	_, allowedTypes, err := h.svc.GetMosqueConfig(c.Request.Context(), mosqueSlug)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "Mosque not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal error"})
		return
	}

	limit := 6
	if l := c.Query("limit"); l != "" {
		if v, err := strconv.Atoi(l); err == nil && v > 0 && v <= 50 {
			limit = v
		}
	}

	types := []string{}
	if t := c.Query("types"); t != "" {
		types = strings.Split(t, ",")
	}

	campaigns, total, err := h.svc.GetCampaigns(c.Request.Context(), repository.CampaignQueryParams{
		AllowedTypes: allowedTypes,
		Types:        types,
		Limit:        limit,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal error"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"campaigns": campaigns,
		"meta":      gin.H{"total": total, "limit": limit},
	})
}

// GetMosqueConfig returns mosque public configuration
func (h *Handler) GetMosqueConfig(c *gin.Context) {
	slug := c.Param("slug")
	mosque, allowedTypes, err := h.svc.GetMosqueConfig(c.Request.Context(), slug)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "Mosque not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal error"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"mosque": gin.H{
			"id":            mosque.ID,
			"name":          mosque.Name,
			"slug":          mosque.Slug,
			"primary_color": mosque.PrimaryColor,
			"logo_url":      mosque.LogoURL,
			"description":   mosque.Description,
			"allowed_types": allowedTypes,
		},
	})
}
