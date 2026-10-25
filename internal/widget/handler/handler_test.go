package handler

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/akordium-id/waqfwise/internal/shared/domain"
	"github.com/akordium-id/waqfwise/internal/widget/repository"
	"github.com/akordium-id/waqfwise/internal/widget/service"
	"github.com/gin-gonic/gin"
)

type stubRepo struct {
	mosque       *domain.Mosque
	allowedTypes []string
	campaigns    []domain.Campaign
	total        int64
	err          error
}

func (s *stubRepo) GetMosqueBySlug(ctx context.Context, slug string) (*domain.Mosque, error) {
	return s.mosque, s.err
}
func (s *stubRepo) GetMosqueCampaignTypes(ctx context.Context, mosqueID string) ([]string, error) {
	return s.allowedTypes, s.err
}
func (s *stubRepo) GetCampaignsForWidget(ctx context.Context, params repository.CampaignQueryParams) ([]domain.Campaign, int64, error) {
	return s.campaigns, s.total, s.err
}

func setupRouter(repo *stubRepo) *gin.Engine {
	gin.SetMode(gin.TestMode)
	svc := service.New(repo, "https://waqfwise.id")
	h := New(svc)
	r := gin.New()
	r.GET("/widget", h.WidgetPage)
	r.GET("/api/v1/widget/campaigns", h.GetCampaigns)
	r.GET("/api/v1/mosques/:slug/config", h.GetMosqueConfig)
	return r
}

func TestGetCampaigns_ReturnsCampaigns(t *testing.T) {
	repo := &stubRepo{
		mosque:       &domain.Mosque{ID: "550e8400-e29b-41d4-a716-446655440000", Slug: "test-mosque", Name: "Test Mosque"},
		allowedTypes: []string{"building"},
		campaigns: []domain.Campaign{
			{
				Title:         "Test Campaign",
				Slug:          "test-campaign",
				Type:          domain.CampaignTypeBuilding,
				GoalAmount:    1000000,
				CurrentAmount: 500000,
				Status:        domain.CampaignStatusActive,
			},
		},
		total: 1,
	}
	router := setupRouter(repo)

	w := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/api/v1/widget/campaigns?mosque_slug=test-mosque&limit=6", nil)
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("status = %d, want 200", w.Code)
	}

	var body map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &body)
	campaigns := body["campaigns"].([]interface{})
	if len(campaigns) != 1 {
		t.Errorf("len(campaigns) = %d, want 1", len(campaigns))
	}
}

func TestGetCampaigns_MissingSlug(t *testing.T) {
	router := setupRouter(&stubRepo{})

	w := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/api/v1/widget/campaigns", nil)
	router.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("status = %d, want 400", w.Code)
	}
}

func TestGetMosqueConfig_ReturnsMosque(t *testing.T) {
	repo := &stubRepo{
		mosque:       &domain.Mosque{ID: "550e8400-e29b-41d4-a716-446655440000", Slug: "test-mosque", Name: "Test Mosque", PrimaryColor: "#1B5E20"},
		allowedTypes: []string{"building", "land"},
	}
	router := setupRouter(repo)

	w := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/api/v1/mosques/test-mosque/config", nil)
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("status = %d, want 200", w.Code)
	}

	var body map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &body)
	mosque := body["mosque"].(map[string]interface{})
	if mosque["name"] != "Test Mosque" {
		t.Errorf("name = %v, want Test Mosque", mosque["name"])
	}
}
