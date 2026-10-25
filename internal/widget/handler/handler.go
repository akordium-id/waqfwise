package handler

import (
	"errors"
	"net/http"
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

	// Referrer validation
	if mosque.Domain != "" {
		referer := c.Request.Referer()
		if referer != "" && !strings.HasPrefix(referer, mosque.Domain) {
			c.String(http.StatusForbidden, "Invalid referrer")
			return
		}
	}

	// Render inline HTML with config
	primaryColor := c.Query("primary_color")
	if primaryColor == "" {
		primaryColor = strings.TrimPrefix(mosque.PrimaryColor, "#")
	}
	logoURL := c.Query("logo_url")
	if logoURL == "" {
		logoURL = mosque.LogoURL
	}
	theme := c.Query("theme")
	if theme == "" {
		theme = "light"
	}

	limit := c.Query("limit")
	if limit == "" {
		limit = "6"
	}
	types := c.Query("types")

	html := renderWidgetHTML(mosque.Name, primaryColor, logoURL, theme, mosqueID, limit, types)
	c.Data(http.StatusOK, "text/html; charset=utf-8", []byte(html))
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

// renderWidgetHTML builds the widget HTML shell
func renderWidgetHTML(mosqueName, primaryColor, logoURL, theme, mosqueID, limit, types string) string {
	return `<!DOCTYPE html>
<html lang="id" data-theme="` + theme + `">
<head>
<meta charset="UTF-8">
<meta name="viewport" content="width=device-width, initial-scale=1.0">
<title>` + mosqueName + ` - Kampanye Wakaf</title>
<style>:root{--waqf-primary:#` + primaryColor + `;--waqf-logo:url(` + logoURL + `)}</style>
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
window.WIDGET_CONFIG = {
  mosqueId: "` + mosqueID + `",
  primaryColor: "#` + primaryColor + `",
  logoUrl: "` + logoURL + `",
  limit: ` + limit + `,
  types: "` + types + `",
  theme: "` + theme + `",
  apiBase: "/api/v1/widget"
};
</script>
<script src="/widget/app.js"></script>
</body>
</html>`
}
