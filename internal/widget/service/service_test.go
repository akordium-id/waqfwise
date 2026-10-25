package service

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/akordium-id/waqfwise/internal/shared/domain"
	"github.com/akordium-id/waqfwise/internal/widget/repository"
)

// mockRepo satisfies RepoInterface for testing.
type mockRepo struct {
	mosque       *domain.Mosque
	allowedTypes []string
	campaigns    []domain.Campaign
	total        int64
	err          error
}

func (m *mockRepo) GetMosqueBySlug(ctx context.Context, slug string) (*domain.Mosque, error) {
	return m.mosque, m.err
}

func (m *mockRepo) GetMosqueCampaignTypes(ctx context.Context, mosqueID string) ([]string, error) {
	return m.allowedTypes, m.err
}

func (m *mockRepo) GetCampaignsForWidget(ctx context.Context, params repository.CampaignQueryParams) ([]domain.Campaign, int64, error) {
	return m.campaigns, m.total, m.err
}

func TestGetCampaigns_TransformsToWidget(t *testing.T) {
	mr := &mockRepo{
		campaigns: []domain.Campaign{
			{
				ID:            1,
				Title:         "Bangun Masjid",
				Slug:          "bangun-masjid",
				ShortDesc:     "Membangun masjid baru",
				Type:          domain.CampaignTypeBuilding,
				GoalAmount:    500000000,
				CurrentAmount: 125000000,
				DonorCount:    342,
				Status:        domain.CampaignStatusActive,
				CreatedAt:     time.Now(),
				UpdatedAt:     time.Now(),
			},
		},
		total: 1,
	}
	svc := New(mr, "https://waqfwise.id")

	campaigns, total, err := svc.GetCampaigns(context.Background(), repository.CampaignQueryParams{Limit: 6})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if total != 1 {
		t.Errorf("total = %d, want 1", total)
	}
	if len(campaigns) != 1 {
		t.Fatalf("len(campaigns) = %d, want 1", len(campaigns))
	}
	c := campaigns[0]
	if c.ID != "bangun-masjid" {
		t.Errorf("ID = %q, want %q", c.ID, "bangun-masjid")
	}
	if c.Title != "Bangun Masjid" {
		t.Errorf("Title = %q, want %q", c.Title, "Bangun Masjid")
	}
	if c.DonateURL != "https://waqfwise.id/donate/bangun-masjid" {
		t.Errorf("DonateURL = %q, want %q", c.DonateURL, "https://waqfwise.id/donate/bangun-masjid")
	}
	if c.Progress != 25.0 {
		t.Errorf("Progress = %v, want 25.0", c.Progress)
	}
	if c.Type != domain.CampaignTypeBuilding {
		t.Errorf("Type = %v, want %v", c.Type, domain.CampaignTypeBuilding)
	}
	if c.GoalAmount != 500000000 {
		t.Errorf("GoalAmount = %d, want 500000000", c.GoalAmount)
	}
	if c.CurrentAmount != 125000000 {
		t.Errorf("CurrentAmount = %d, want 125000000", c.CurrentAmount)
	}
	if c.DonorCount != 342 {
		t.Errorf("DonorCount = %d, want 342", c.DonorCount)
	}
}

func TestGetCampaigns_ZeroGoalProgress(t *testing.T) {
	mr := &mockRepo{
		campaigns: []domain.Campaign{
			{
				ID:            2,
				Title:         "Zero Goal",
				Slug:          "zero-goal",
				GoalAmount:    0,
				CurrentAmount: 100000,
				Status:        domain.CampaignStatusActive,
			},
		},
		total: 1,
	}
	svc := New(mr, "https://waqfwise.id")

	campaigns, _, err := svc.GetCampaigns(context.Background(), repository.CampaignQueryParams{Limit: 6})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if campaigns[0].Progress != 0 {
		t.Errorf("Progress = %v, want 0 for zero goal", campaigns[0].Progress)
	}
}

func TestGetCampaigns_PropagatesRepoError(t *testing.T) {
	expected := errors.New("db down")
	mr := &mockRepo{err: expected}
	svc := New(mr, "https://waqfwise.id")

	_, _, err := svc.GetCampaigns(context.Background(), repository.CampaignQueryParams{Limit: 6})
	if !errors.Is(err, expected) {
		t.Errorf("err = %v, want wrapped %v", err, expected)
	}
}

func TestGetMosqueConfig_ReturnsMosqueAndTypes(t *testing.T) {
	mosque := &domain.Mosque{
		ID:   "mosque-1",
		Name: "Masjid Al-Ikhlas",
		Slug: "masjid-al-ikhlas",
	}
	mr := &mockRepo{
		mosque:       mosque,
		allowedTypes: []string{"building", "cash"},
	}
	svc := New(mr, "https://waqfwise.id")

	m, types, err := svc.GetMosqueConfig(context.Background(), "masjid-al-ikhlas")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if m.Slug != "masjid-al-ikhlas" {
		t.Errorf("slug = %q, want %q", m.Slug, "masjid-al-ikhlas")
	}
	if len(types) != 2 {
		t.Errorf("len(types) = %d, want 2", len(types))
	}
}

func TestGetMosqueConfig_PropagatesRepoError(t *testing.T) {
	expected := errors.New("db down")
	mr := &mockRepo{err: expected}
	svc := New(mr, "https://waqfwise.id")

	_, _, err := svc.GetMosqueConfig(context.Background(), "slug")
	if !errors.Is(err, expected) {
		t.Errorf("err = %v, want %v", err, expected)
	}
}
