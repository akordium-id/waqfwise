package service

import (
	"context"
	"encoding/json"
	"strings"

	"github.com/akordium-id/waqfwise/internal/shared/domain"
)

// FraudDetector handles fraud detection
type FraudDetector struct {
	// Add third-party fraud detection service here
}

// NewFraudDetector creates a new fraud detector
func NewFraudDetector() *FraudDetector {
	return &FraudDetector{}
}

// CheckTransaction performs fraud detection on transaction
func (f *FraudDetector) CheckTransaction(ctx context.Context, donation *domain.Donation, ipAddress, deviceID string) (*domain.FraudCheck, error) {
	riskScore := 0
	flags := make([]string, 0)

	// Check 1: Amount threshold
	if donation.Amount > 100000000 { // > 100 juta
		riskScore += 30
		flags = append(flags, "high_amount")
	}

	// Check 2: Suspicious email patterns
	if strings.Contains(donation.DonorEmail, "temp") || strings.Contains(donation.DonorEmail, "disposable") {
		riskScore += 20
		flags = append(flags, "suspicious_email")
	}

	// Check 3: Multiple small transactions (velocity check)
	// TODO: Implement velocity checking from database
	// if velocityCount > threshold {
	//     riskScore += 25
	//     flags = append(flags, "high_velocity")
	// }

	// Check 4: IP blacklist
	// TODO: Check IP against blacklist database

	// Determine risk level
	var riskLevel string
	var isBlocked bool
	var reason string

	if riskScore >= 70 {
		riskLevel = "high"
		isBlocked = true
		reason = "High risk transaction blocked"
	} else if riskScore >= 40 {
		riskLevel = "medium"
		isBlocked = false
		reason = "Medium risk - requires manual review"
	} else {
		riskLevel = "low"
		isBlocked = false
		reason = ""
	}

	flagsJSON, _ := json.Marshal(flags)

	return &domain.FraudCheck{
		DonationID: donation.ID,
		RiskScore:  riskScore,
		RiskLevel:  riskLevel,
		Flags:      string(flagsJSON),
		IsBlocked:  isBlocked,
		Reason:     reason,
		IPAddress:  ipAddress,
		DeviceID:   deviceID,
	}, nil
}
