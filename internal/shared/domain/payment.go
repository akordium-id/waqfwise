package domain

import (
	"time"
)

// PaymentStatus represents the status of a payment
type PaymentStatus string

const (
	PaymentStatusPending   PaymentStatus = "pending"
	PaymentStatusProcessing PaymentStatus = "processing"
	PaymentStatusSuccess   PaymentStatus = "success"
	PaymentStatusFailed    PaymentStatus = "failed"
	PaymentStatusRefunded  PaymentStatus = "refunded"
	PaymentStatusCancelled PaymentStatus = "cancelled"
)

// PaymentMethod represents the payment method used
type PaymentMethod string

const (
	PaymentMethodCreditCard PaymentMethod = "credit_card"
	PaymentMethodBankTransfer PaymentMethod = "bank_transfer"
	PaymentMethodEWallet    PaymentMethod = "e_wallet"
	PaymentMethodQRIS       PaymentMethod = "qris"
	PaymentMethodVA         PaymentMethod = "virtual_account"
)

// PaymentGateway represents the payment gateway provider
type PaymentGateway string

const (
	PaymentGatewayMidtrans PaymentGateway = "midtrans"
	PaymentGatewayXendit   PaymentGateway = "xendit"
	PaymentGatewayManual   PaymentGateway = "manual"
)

// Donation represents a donation transaction
type Donation struct {
	ID              int64          `json:"id" db:"id"`
	CampaignID      int64          `json:"campaign_id" db:"campaign_id"`
	UserID          int64          `json:"user_id" db:"user_id"`
	Amount          int64          `json:"amount" db:"amount"`
	Status          PaymentStatus  `json:"status" db:"status"`
	PaymentMethod   PaymentMethod  `json:"payment_method" db:"payment_method"`
	PaymentGateway  PaymentGateway `json:"payment_gateway" db:"payment_gateway"`
	TransactionID   string         `json:"transaction_id" db:"transaction_id"`
	GatewayRef      string         `json:"gateway_ref,omitempty" db:"gateway_ref"`
	IsAnonymous     bool           `json:"is_anonymous" db:"is_anonymous"`
	DonorName       string         `json:"donor_name,omitempty" db:"donor_name"`
	DonorEmail      string         `json:"donor_email,omitempty" db:"donor_email"`
	Message         string         `json:"message,omitempty" db:"message"`
	IsRecurring     bool           `json:"is_recurring" db:"is_recurring"`
	RecurringPeriod string         `json:"recurring_period,omitempty" db:"recurring_period"`
	ReceiptURL      string         `json:"receipt_url,omitempty" db:"receipt_url"`
	PaidAt          *time.Time     `json:"paid_at,omitempty" db:"paid_at"`
	CreatedAt       time.Time      `json:"created_at" db:"created_at"`
	UpdatedAt       time.Time      `json:"updated_at" db:"updated_at"`
}

// PaymentLog represents detailed payment logs
type PaymentLog struct {
	ID           int64          `json:"id" db:"id"`
	DonationID   int64          `json:"donation_id" db:"donation_id"`
	Status       PaymentStatus  `json:"status" db:"status"`
	Gateway      PaymentGateway `json:"gateway" db:"gateway"`
	RequestData  string         `json:"request_data" db:"request_data"`
	ResponseData string         `json:"response_data" db:"response_data"`
	ErrorMessage string         `json:"error_message,omitempty" db:"error_message"`
	IPAddress    string         `json:"ip_address,omitempty" db:"ip_address"`
	UserAgent    string         `json:"user_agent,omitempty" db:"user_agent"`
	CreatedAt    time.Time      `json:"created_at" db:"created_at"`
}

// Ledger represents double-entry bookkeeping ledger
type Ledger struct {
	ID             int64     `json:"id" db:"id"`
	DonationID     int64     `json:"donation_id" db:"donation_id"`
	AccountType    string    `json:"account_type" db:"account_type"` // debit/credit
	AccountName    string    `json:"account_name" db:"account_name"` // campaign_fund, gateway_fee, net_revenue
	Amount         int64     `json:"amount" db:"amount"`
	BalanceBefore  int64     `json:"balance_before" db:"balance_before"`
	BalanceAfter   int64     `json:"balance_after" db:"balance_after"`
	Description    string    `json:"description" db:"description"`
	CreatedAt      time.Time `json:"created_at" db:"created_at"`
}

// FraudCheck represents fraud detection results
type FraudCheck struct {
	ID          int64     `json:"id" db:"id"`
	DonationID  int64     `json:"donation_id" db:"donation_id"`
	RiskScore   int       `json:"risk_score" db:"risk_score"`       // 0-100
	RiskLevel   string    `json:"risk_level" db:"risk_level"`       // low/medium/high
	Flags       string    `json:"flags" db:"flags"`                 // JSON array of flags
	IsBlocked   bool      `json:"is_blocked" db:"is_blocked"`
	Reason      string    `json:"reason,omitempty" db:"reason"`
	IPAddress   string    `json:"ip_address" db:"ip_address"`
	DeviceID    string    `json:"device_id,omitempty" db:"device_id"`
	CreatedAt   time.Time `json:"created_at" db:"created_at"`
}

// IsPaid checks if donation is paid
func (d *Donation) IsPaid() bool {
	return d.Status == PaymentStatusSuccess
}

// IsPending checks if donation is pending
func (d *Donation) IsPending() bool {
	return d.Status == PaymentStatusPending || d.Status == PaymentStatusProcessing
}
