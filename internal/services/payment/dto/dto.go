package dto

import "github.com/akordium-id/waqfwise/internal/shared/domain"

// CreateDonationRequest represents donation creation request
type CreateDonationRequest struct {
	CampaignID      int64                   `json:"campaign_id"`
	Amount          int64                   `json:"amount"`
	PaymentMethod   domain.PaymentMethod    `json:"payment_method"`
	PaymentGateway  domain.PaymentGateway   `json:"payment_gateway"`
	IsAnonymous     bool                    `json:"is_anonymous"`
	DonorName       string                  `json:"donor_name,omitempty"`
	DonorEmail      string                  `json:"donor_email,omitempty"`
	Message         string                  `json:"message,omitempty"`
	IsRecurring     bool                    `json:"is_recurring"`
	RecurringPeriod string                  `json:"recurring_period,omitempty"` // monthly, yearly
}

// PaymentCallbackRequest represents payment gateway callback
type PaymentCallbackRequest struct {
	Gateway       domain.PaymentGateway `json:"gateway"`
	TransactionID string                `json:"transaction_id"`
	Status        string                `json:"status"`
	RawData       map[string]interface{} `json:"raw_data"`
}

// DonationResponse represents donation response
type DonationResponse struct {
	ID              int64                   `json:"id"`
	CampaignID      int64                   `json:"campaign_id"`
	UserID          int64                   `json:"user_id"`
	Amount          int64                   `json:"amount"`
	Status          domain.PaymentStatus    `json:"status"`
	PaymentMethod   domain.PaymentMethod    `json:"payment_method"`
	PaymentGateway  domain.PaymentGateway   `json:"payment_gateway"`
	TransactionID   string                  `json:"transaction_id"`
	PaymentURL      string                  `json:"payment_url,omitempty"`
	IsAnonymous     bool                    `json:"is_anonymous"`
	Message         string                  `json:"message,omitempty"`
	ReceiptURL      string                  `json:"receipt_url,omitempty"`
	CreatedAt       string                  `json:"created_at"`
}

// LedgerResponse represents ledger entry response
type LedgerResponse struct {
	ID            int64  `json:"id"`
	DonationID    int64  `json:"donation_id"`
	AccountType   string `json:"account_type"`
	AccountName   string `json:"account_name"`
	Amount        int64  `json:"amount"`
	BalanceBefore int64  `json:"balance_before"`
	BalanceAfter  int64  `json:"balance_after"`
	Description   string `json:"description"`
	CreatedAt     string `json:"created_at"`
}

// FraudCheckResponse represents fraud check result
type FraudCheckResponse struct {
	IsBlocked bool   `json:"is_blocked"`
	RiskLevel string `json:"risk_level"`
	RiskScore int    `json:"risk_score"`
	Reason    string `json:"reason,omitempty"`
}

// FromDomain converts domain.Donation to DonationResponse
func FromDomain(donation *domain.Donation) *DonationResponse {
	return &DonationResponse{
		ID:             donation.ID,
		CampaignID:     donation.CampaignID,
		UserID:         donation.UserID,
		Amount:         donation.Amount,
		Status:         donation.Status,
		PaymentMethod:  donation.PaymentMethod,
		PaymentGateway: donation.PaymentGateway,
		TransactionID:  donation.TransactionID,
		IsAnonymous:    donation.IsAnonymous,
		Message:        donation.Message,
		ReceiptURL:     donation.ReceiptURL,
		CreatedAt:      donation.CreatedAt.Format("2006-01-02T15:04:05Z"),
	}
}
