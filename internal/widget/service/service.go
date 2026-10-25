package service

import (
	"context"

	"github.com/akordium-id/waqfwise/internal/shared/domain"
	"github.com/akordium-id/waqfwise/internal/widget/repository"
)

// RepoInterface abstracts the repository for testability
type RepoInterface interface {
	GetMosqueBySlug(ctx context.Context, slug string) (*domain.Mosque, error)
	GetMosqueCampaignTypes(ctx context.Context, mosqueID string) ([]string, error)
	GetCampaignsForWidget(ctx context.Context, params repository.CampaignQueryParams) ([]domain.Campaign, int64, error)
}

// Service provides business logic for the widget
type Service struct {
	repo    RepoInterface
	baseURL string
}

// New creates a new widget service
func New(repo RepoInterface, baseURL string) *Service {
	return &Service{repo: repo, baseURL: baseURL}
}

// GetMosqueConfig returns the mosque and its allowed campaign types
func (s *Service) GetMosqueConfig(ctx context.Context, slug string) (*domain.Mosque, []string, error) {
	mosque, err := s.repo.GetMosqueBySlug(ctx, slug)
	if err != nil {
		return nil, nil, err
	}
	types, err := s.repo.GetMosqueCampaignTypes(ctx, mosque.ID)
	if err != nil {
		return nil, nil, err
	}
	return mosque, types, nil
}

// GetCampaigns returns widget-formatted campaigns
func (s *Service) GetCampaigns(ctx context.Context, params repository.CampaignQueryParams) ([]domain.WidgetCampaign, int64, error) {
	campaigns, total, err := s.repo.GetCampaignsForWidget(ctx, params)
	if err != nil {
		return nil, 0, err
	}

	widgetCampaigns := make([]domain.WidgetCampaign, len(campaigns))
	for i, c := range campaigns {
		wc := domain.WidgetCampaign{
			ID:            c.Slug,
			Title:         c.Title,
			Slug:          c.Slug,
			ShortDesc:     c.ShortDesc,
			Type:          c.Type,
			GoalAmount:    c.GoalAmount,
			CurrentAmount: c.CurrentAmount,
			DonorCount:    c.DonorCount,
			ImageURL:      c.ImageURL,
			IsUrgent:      c.IsUrgent,
			IsFeatured:    c.IsFeatured,
			DonateURL:     s.baseURL + "/donate/" + c.Slug,
		}
		if wc.GoalAmount == 0 {
			wc.Progress = 0
		} else {
			wc.Progress = (float64(wc.CurrentAmount) / float64(wc.GoalAmount)) * 100
		}
		widgetCampaigns[i] = wc
	}
	return widgetCampaigns, total, nil
}
