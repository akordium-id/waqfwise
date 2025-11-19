package payment

import (
	"context"
	"time"
)

// PaymentStatus represents the status of a payment
type PaymentStatus string

const (
	StatusPending   PaymentStatus = "pending"
	StatusSuccess   PaymentStatus = "success"
	StatusFailed    PaymentStatus = "failed"
	StatusCancelled PaymentStatus = "cancelled"
	StatusExpired   PaymentStatus = "expired"
)

// PaymentMethod represents a payment method
type PaymentMethod string

const (
	MethodCreditCard   PaymentMethod = "credit_card"
	MethodBankTransfer PaymentMethod = "bank_transfer"
	MethodEWallet      PaymentMethod = "e_wallet"
	MethodQRIS         PaymentMethod = "qris"
	MethodVA           PaymentMethod = "virtual_account"
)

// PaymentRequest represents a payment request
type PaymentRequest struct {
	OrderID       string
	Amount        int64
	Currency      string
	CustomerName  string
	CustomerEmail string
	CustomerPhone string
	Description   string
	Items         []PaymentItem
	Metadata      map[string]interface{}
}

// PaymentItem represents an item in a payment
type PaymentItem struct {
	ID       string
	Name     string
	Price    int64
	Quantity int
}

// PaymentResponse represents a payment response
type PaymentResponse struct {
	TransactionID string
	OrderID       string
	Status        PaymentStatus
	Amount        int64
	PaymentURL    string
	QRCode        string
	VANumber      string
	ExpiredAt     time.Time
	PaidAt        *time.Time
	Metadata      map[string]interface{}
}

// PaymentNotification represents a payment notification/webhook
type PaymentNotification struct {
	TransactionID string
	OrderID       string
	Status        PaymentStatus
	Amount        int64
	PaidAt        *time.Time
	Metadata      map[string]interface{}
}

// PaymentGateway defines the interface for payment gateways
type PaymentGateway interface {
	// CreateTransaction creates a new payment transaction
	CreateTransaction(ctx context.Context, req *PaymentRequest) (*PaymentResponse, error)

	// GetTransaction retrieves a transaction by ID
	GetTransaction(ctx context.Context, transactionID string) (*PaymentResponse, error)

	// CancelTransaction cancels a transaction
	CancelTransaction(ctx context.Context, transactionID string) error

	// VerifyNotification verifies a payment notification
	VerifyNotification(ctx context.Context, payload map[string]interface{}) (*PaymentNotification, error)

	// GetName returns the name of the payment gateway
	GetName() string
}
