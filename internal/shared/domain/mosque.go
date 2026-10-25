package domain

import (
	"time"
)

// Mosque represents a mosque with widget configuration
type Mosque struct {
	ID           string    `json:"id" gorm:"type:uuid;primaryKey;default:uuid_generate_v4()"`
	Name         string    `json:"name" gorm:"column:name;size:255;not null"`
	Slug         string    `json:"slug" gorm:"column:slug;size:100;uniqueIndex;not null"`
	Domain       string    `json:"domain,omitempty" gorm:"column:domain;size:255"`
	PrimaryColor string    `json:"primary_color" gorm:"column:primary_color;size:7;default:#1B5E20"`
	LogoURL      string    `json:"logo_url,omitempty" gorm:"column:logo_url;type:text"`
	Description  string    `json:"description,omitempty" gorm:"column:description;type:text"`
	Address      string    `json:"address,omitempty" gorm:"column:address;type:text"`
	Phone        string    `json:"phone,omitempty" gorm:"column:phone;size:50"`
	Email        string    `json:"email,omitempty" gorm:"column:email;size:255"`
	IsActive     bool      `json:"is_active" gorm:"column:is_active;default:true"`
	CreatedAt    time.Time `json:"created_at" gorm:"column:created_at;autoCreateTime"`
	UpdatedAt    time.Time `json:"updated_at" gorm:"column:updated_at;autoUpdateTime"`
}

func (Mosque) TableName() string { return "mosques" }

// MosqueCampaignType maps allowed campaign types per mosque
type MosqueCampaignType struct {
	MosqueID     string `json:"mosque_id" gorm:"column:mosque_id;type:uuid;primaryKey"`
	CampaignType string `json:"campaign_type" gorm:"column:campaign_type;size:50;primaryKey"`
}

func (MosqueCampaignType) TableName() string { return "mosque_campaign_types" }

// WidgetCampaign is the public-facing campaign DTO (no internal fields)
type WidgetCampaign struct {
	ID            string       `json:"id"`
	Title         string       `json:"title"`
	Slug          string       `json:"slug"`
	ShortDesc     string       `json:"short_desc"`
	Type          CampaignType `json:"type"`
	GoalAmount    int64        `json:"goal_amount"`
	CurrentAmount int64        `json:"current_amount"`
	DonorCount    int          `json:"donor_count"`
	ImageURL      string       `json:"image_url,omitempty"`
	Progress      float64      `json:"progress"`
	IsUrgent      bool         `json:"is_urgent"`
	IsFeatured    bool         `json:"is_featured"`
	DonateURL     string       `json:"donate_url"`
}
