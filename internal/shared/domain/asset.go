package domain

import (
	"time"
)

// AssetType represents the type of wakaf asset
type AssetType string

const (
	AssetTypeLand     AssetType = "land"
	AssetTypeBuilding AssetType = "building"
	AssetTypeVehicle  AssetType = "vehicle"
	AssetTypeEquipment AssetType = "equipment"
	AssetTypeOther    AssetType = "other"
)

// AssetStatus represents the status of an asset
type AssetStatus string

const (
	AssetStatusActive     AssetStatus = "active"
	AssetStatusInactive   AssetStatus = "inactive"
	AssetStatusUnderMaint AssetStatus = "under_maintenance"
	AssetStatusDisposed   AssetStatus = "disposed"
)

// Asset represents a wakaf asset
type Asset struct {
	ID              int64       `json:"id" db:"id"`
	CampaignID      int64       `json:"campaign_id" db:"campaign_id"`
	Type            AssetType   `json:"type" db:"type"`
	Status          AssetStatus `json:"status" db:"status"`
	Name            string      `json:"name" db:"name"`
	Description     string      `json:"description" db:"description"`
	Location        string      `json:"location" db:"location"`
	Latitude        *float64    `json:"latitude,omitempty" db:"latitude"`
	Longitude       *float64    `json:"longitude,omitempty" db:"longitude"`
	Area            *float64    `json:"area,omitempty" db:"area"`           // in sqm
	AreaUnit        string      `json:"area_unit,omitempty" db:"area_unit"` // sqm, hectare
	PurchaseValue   int64       `json:"purchase_value" db:"purchase_value"`
	CurrentValue    int64       `json:"current_value" db:"current_value"`
	LastValuationAt *time.Time  `json:"last_valuation_at,omitempty" db:"last_valuation_at"`
	AcquisitionDate time.Time   `json:"acquisition_date" db:"acquisition_date"`
	TenantID        *int64      `json:"tenant_id,omitempty" db:"tenant_id"`
	CreatedAt       time.Time   `json:"created_at" db:"created_at"`
	UpdatedAt       time.Time   `json:"updated_at" db:"updated_at"`
}

// AssetDocument represents legal documents for assets
type AssetDocument struct {
	ID           int64     `json:"id" db:"id"`
	AssetID      int64     `json:"asset_id" db:"asset_id"`
	DocumentType string    `json:"document_type" db:"document_type"` // certificate, deed, permit
	DocumentName string    `json:"document_name" db:"document_name"`
	FileURL      string    `json:"file_url" db:"file_url"`
	FileSize     int64     `json:"file_size" db:"file_size"`
	MimeType     string    `json:"mime_type" db:"mime_type"`
	IssueDate    *time.Time `json:"issue_date,omitempty" db:"issue_date"`
	ExpiryDate   *time.Time `json:"expiry_date,omitempty" db:"expiry_date"`
	Notes        string    `json:"notes,omitempty" db:"notes"`
	UploadedBy   int64     `json:"uploaded_by" db:"uploaded_by"`
	CreatedAt    time.Time `json:"created_at" db:"created_at"`
	UpdatedAt    time.Time `json:"updated_at" db:"updated_at"`
}

// AssetValuation represents asset valuation history
type AssetValuation struct {
	ID             int64     `json:"id" db:"id"`
	AssetID        int64     `json:"asset_id" db:"asset_id"`
	ValuationValue int64     `json:"valuation_value" db:"valuation_value"`
	ValuationDate  time.Time `json:"valuation_date" db:"valuation_date"`
	ValuedBy       string    `json:"valued_by" db:"valued_by"` // appraiser name
	Method         string    `json:"method" db:"method"`        // market_value, income_approach, etc
	Notes          string    `json:"notes,omitempty" db:"notes"`
	DocumentURL    string    `json:"document_url,omitempty" db:"document_url"`
	CreatedAt      time.Time `json:"created_at" db:"created_at"`
}

// GeospatialData represents additional geospatial information
type GeospatialData struct {
	AssetID     int64     `json:"asset_id" db:"asset_id"`
	Geometry    string    `json:"geometry" db:"geometry"`         // PostGIS geometry
	Boundary    string    `json:"boundary,omitempty" db:"boundary"` // GeoJSON polygon
	LandUse     string    `json:"land_use,omitempty" db:"land_use"`
	ZoningType  string    `json:"zoning_type,omitempty" db:"zoning_type"`
	Elevation   *float64  `json:"elevation,omitempty" db:"elevation"`
	CreatedAt   time.Time `json:"created_at" db:"created_at"`
	UpdatedAt   time.Time `json:"updated_at" db:"updated_at"`
}

// IsLand checks if asset is land type
func (a *Asset) IsLand() bool {
	return a.Type == AssetTypeLand
}

// HasGeolocation checks if asset has geolocation data
func (a *Asset) HasGeolocation() bool {
	return a.Latitude != nil && a.Longitude != nil
}
