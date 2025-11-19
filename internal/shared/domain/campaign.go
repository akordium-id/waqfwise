package domain

import (
	"time"
)

// CampaignStatus represents the status of a campaign
type CampaignStatus string

const (
	CampaignStatusDraft     CampaignStatus = "draft"
	CampaignStatusActive    CampaignStatus = "active"
	CampaignStatusCompleted CampaignStatus = "completed"
	CampaignStatusCancelled CampaignStatus = "cancelled"
	CampaignStatusPaused    CampaignStatus = "paused"
)

// CampaignType represents the type of wakaf
type CampaignType string

const (
	CampaignTypeLand       CampaignType = "land"
	CampaignTypeBuilding   CampaignType = "building"
	CampaignTypeCash       CampaignType = "cash"
	CampaignTypeEducation  CampaignType = "education"
	CampaignTypeHealthcare CampaignType = "healthcare"
	CampaignTypeGeneral    CampaignType = "general"
)

// Campaign represents a wakaf campaign
type Campaign struct {
	ID              int64          `json:"id" db:"id"`
	Title           string         `json:"title" db:"title"`
	Slug            string         `json:"slug" db:"slug"`
	Description     string         `json:"description" db:"description"`
	ShortDesc       string         `json:"short_desc" db:"short_desc"`
	Type            CampaignType   `json:"type" db:"type"`
	Status          CampaignStatus `json:"status" db:"status"`
	GoalAmount      int64          `json:"goal_amount" db:"goal_amount"`
	CurrentAmount   int64          `json:"current_amount" db:"current_amount"`
	DonorCount      int            `json:"donor_count" db:"donor_count"`
	NazirID         int64          `json:"nazir_id" db:"nazir_id"`
	ImageURL        string         `json:"image_url,omitempty" db:"image_url"`
	VideoURL        string         `json:"video_url,omitempty" db:"video_url"`
	Location        string         `json:"location,omitempty" db:"location"`
	StartDate       time.Time      `json:"start_date" db:"start_date"`
	EndDate         *time.Time     `json:"end_date,omitempty" db:"end_date"`
	IsEndless       bool           `json:"is_endless" db:"is_endless"`
	IsFeatured      bool           `json:"is_featured" db:"is_featured"`
	IsUrgent        bool           `json:"is_urgent" db:"is_urgent"`
	TenantID        *int64         `json:"tenant_id,omitempty" db:"tenant_id"`
	CreatedAt       time.Time      `json:"created_at" db:"created_at"`
	UpdatedAt       time.Time      `json:"updated_at" db:"updated_at"`
}

// Milestone represents campaign milestones
type Milestone struct {
	ID          int64     `json:"id" db:"id"`
	CampaignID  int64     `json:"campaign_id" db:"campaign_id"`
	Title       string    `json:"title" db:"title"`
	Description string    `json:"description" db:"description"`
	TargetAmount int64    `json:"target_amount" db:"target_amount"`
	IsCompleted bool      `json:"is_completed" db:"is_completed"`
	CompletedAt *time.Time `json:"completed_at,omitempty" db:"completed_at"`
	CreatedAt   time.Time `json:"created_at" db:"created_at"`
	UpdatedAt   time.Time `json:"updated_at" db:"updated_at"`
}

// Progress calculates campaign progress percentage
func (c *Campaign) Progress() float64 {
	if c.GoalAmount == 0 {
		return 0
	}
	return (float64(c.CurrentAmount) / float64(c.GoalAmount)) * 100
}

// IsActive checks if campaign is active
func (c *Campaign) IsActive() bool {
	return c.Status == CampaignStatusActive
}

// IsCompleted checks if campaign has reached its goal
func (c *Campaign) IsCompleted() bool {
	return c.CurrentAmount >= c.GoalAmount
}
