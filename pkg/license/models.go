package license

import "time"

// License represents a WaqfWise license
type License struct {
	ID           string    `json:"id"`
	CustomerID   string    `json:"customer_id"`
	CustomerName string    `json:"customer_name"`
	Features     []string  `json:"features"`
	ExpiresAt    time.Time `json:"expires_at"`
	IssuedAt     time.Time `json:"issued_at"`
	MaxUsers     int       `json:"max_users"`
	Signature    string    `json:"signature"`
}

// Feature represents an enterprise feature
type Feature string

const (
	// FeatureMultiTenancy enables multi-tenant architecture
	FeatureMultiTenancy Feature = "multitenancy"

	// FeatureAnalytics enables advanced analytics engine
	FeatureAnalytics Feature = "analytics"

	// FeatureWhiteLabel enables white-labeling
	FeatureWhiteLabel Feature = "whitelabel"

	// FeatureIntegration enables third-party integrations
	FeatureIntegration Feature = "integration"

	// FeatureAdvancedPayment enables advanced payment features
	FeatureAdvancedPayment Feature = "advanced_payment"

	// FeatureGeospatial enables geospatial features
	FeatureGeospatial Feature = "geospatial"
)

// HasFeature checks if license has a specific feature
func (l *License) HasFeature(feature Feature) bool {
	for _, f := range l.Features {
		if f == string(feature) {
			return true
		}
	}
	return false
}

// IsValid checks if license is currently valid
func (l *License) IsValid() bool {
	return time.Now().Before(l.ExpiresAt)
}
