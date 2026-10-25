package repository

import (
	"context"

	"github.com/akordium-id/waqfwise/internal/shared/domain"
	"gorm.io/gorm"
)

// CampaignQueryParams holds filtering options for widget campaigns
type CampaignQueryParams struct {
	AllowedTypes []string
	Types        []string
	Limit        int
}

// Repository provides data access for widget queries
type Repository struct {
	db *gorm.DB
}

// New creates a new widget repository
func New(db *gorm.DB) *Repository {
	return &Repository{db: db}
}

// GetMosqueBySlug returns an active mosque by its slug
func (r *Repository) GetMosqueBySlug(ctx context.Context, slug string) (*domain.Mosque, error) {
	var m domain.Mosque
	err := r.db.WithContext(ctx).
		Where("slug = ? AND is_active = ?", slug, true).
		First(&m).Error
	if err != nil {
		return nil, err
	}
	return &m, nil
}

// GetMosqueCampaignTypes returns allowed campaign types for a mosque
func (r *Repository) GetMosqueCampaignTypes(ctx context.Context, mosqueID string) ([]string, error) {
	var types []string
	err := r.db.WithContext(ctx).
		Model(&domain.MosqueCampaignType{}).
		Where("mosque_id = ?", mosqueID).
		Pluck("campaign_type", &types).Error
	return types, err
}

// GetCampaignsForWidget returns active campaigns matching filters
func (r *Repository) GetCampaignsForWidget(ctx context.Context, params CampaignQueryParams) ([]domain.Campaign, int64, error) {
	query := r.db.WithContext(ctx).
		Model(&domain.Campaign{}).
		Where("status = ? AND deleted_at IS NULL", "active")

	if len(params.AllowedTypes) > 0 {
		query = query.Where("type IN ?", params.AllowedTypes)
	}
	if len(params.Types) > 0 {
		query = query.Where("type IN ?", params.Types)
	}

	var total int64
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	var campaigns []domain.Campaign
	err := query.
		Order("is_featured DESC, is_urgent DESC, created_at DESC").
		Limit(params.Limit).
		Find(&campaigns).Error

	return campaigns, total, err
}
