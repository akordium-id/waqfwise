package handler

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/akordium-id/waqfwise/internal/shared/domain"
	"github.com/akordium-id/waqfwise/internal/widget/repository"
	"github.com/akordium-id/waqfwise/internal/widget/service"
	"github.com/gin-gonic/gin"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func setupTestDB(t *testing.T) *gorm.DB {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Fatalf("failed to open test db: %v", err)
	}
	db.Exec(`CREATE TABLE mosques (id TEXT PRIMARY KEY, name TEXT NOT NULL, slug TEXT UNIQUE NOT NULL, domain TEXT, primary_color TEXT DEFAULT '#1B5E20', logo_url TEXT, description TEXT, address TEXT, phone TEXT, email TEXT, is_active INTEGER NOT NULL DEFAULT 1, created_at DATETIME, updated_at DATETIME)`)
	db.Exec(`CREATE TABLE mosque_campaign_types (mosque_id TEXT NOT NULL, campaign_type TEXT NOT NULL, PRIMARY KEY (mosque_id, campaign_type))`)
	db.Exec(`CREATE TABLE campaigns (id INTEGER PRIMARY KEY AUTOINCREMENT, title TEXT NOT NULL, slug TEXT NOT NULL, description TEXT, short_desc TEXT, type TEXT NOT NULL, status TEXT NOT NULL, goal_amount INTEGER NOT NULL, current_amount INTEGER NOT NULL, donor_count INTEGER NOT NULL, nazir_id INTEGER NOT NULL, image_url TEXT, video_url TEXT, location TEXT, start_date DATETIME NOT NULL, end_date DATETIME, is_endless INTEGER NOT NULL, is_featured INTEGER NOT NULL, is_urgent INTEGER NOT NULL, tenant_id INTEGER, deleted_at DATETIME, created_at DATETIME NOT NULL, updated_at DATETIME NOT NULL)`)
	return db
}

func TestIntegration_FullWidgetFlow(t *testing.T) {
	db := setupTestDB(t)

	// Seed mosque
	mosque := domain.Mosque{
		ID:           "550e8400-e29b-41d4-a716-446655440000",
		Name:         "Masjid Al-Ikhlas",
		Slug:         "masjid-al-ikhlas",
		PrimaryColor: "#1B5E20",
		IsActive:     true,
	}
	db.Create(&mosque)
	db.Create(&domain.MosqueCampaignType{
		MosqueID:     mosque.ID,
		CampaignType: "building",
	})

	// Seed campaigns
	db.Create(&domain.Campaign{
		Title:         "Bangun Masjid",
		Slug:          "bangun-masjid",
		Type:          domain.CampaignTypeBuilding,
		GoalAmount:    500000000,
		CurrentAmount: 125000000,
		Status:        domain.CampaignStatusActive,
	})
	db.Create(&domain.Campaign{
		Title:  "Draft Campaign",
		Slug:   "draft-campaign",
		Type:   domain.CampaignTypeBuilding,
		Status: domain.CampaignStatusDraft, // Should NOT appear
	})

	repo := repository.New(db)
	svc := service.New(repo, "https://waqfwise.id")
	h := New(svc)

	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.GET("/widget", h.WidgetPage)
	r.GET("/api/v1/widget/campaigns", h.GetCampaigns)

	// Test 1: Widget page renders
	w := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/widget?mosque_id=masjid-al-ikhlas", nil)
	r.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Fatalf("widget page status = %d, want 200", w.Code)
	}

	// Test 2: Campaign API returns only active
	w2 := httptest.NewRecorder()
	req2 := httptest.NewRequest("GET", "/api/v1/widget/campaigns?mosque_slug=masjid-al-ikhlas", nil)
	r.ServeHTTP(w2, req2)
	if w2.Code != http.StatusOK {
		t.Fatalf("campaigns status = %d, want 200", w2.Code)
	}

	var body map[string]interface{}
	json.Unmarshal(w2.Body.Bytes(), &body)
	campaigns := body["campaigns"].([]interface{})
	if len(campaigns) != 1 {
		t.Errorf("len(campaigns) = %d, want 1 (only active)", len(campaigns))
	}

	// Test 3: donate_url is correct
	first := campaigns[0].(map[string]interface{})
	if first["donate_url"] != "https://waqfwise.id/donate/bangun-masjid" {
		t.Errorf("donate_url = %v, want https://waqfwise.id/donate/bangun-masjid", first["donate_url"])
	}
}
