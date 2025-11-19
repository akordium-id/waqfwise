package payment

import (
	"context"
	"fmt"
	"time"

	"github.com/akordium-id/waqfwise/pkg/config"
	"github.com/xendit/xendit-go/v4"
	"go.uber.org/zap"
)

// XenditGateway implements PaymentGateway for Xendit
type XenditGateway struct {
	client *xendit.APIClient
	config *config.XenditConfig
	logger *zap.Logger
}

// NewXenditGateway creates a new Xendit payment gateway
func NewXenditGateway(cfg *config.XenditConfig, logger *zap.Logger) *XenditGateway {
	// Create Xendit client configuration
	xenditConfig := xendit.NewConfiguration()
	xenditConfig.SetHTTPClient(&xendit.DefaultHTTPClient{
		SecretKey: cfg.SecretKey,
	})

	client := xendit.NewAPIClient(xenditConfig)

	return &XenditGateway{
		client: client,
		config: cfg,
		logger: logger,
	}
}

// CreateTransaction creates a new payment transaction
func (x *XenditGateway) CreateTransaction(ctx context.Context, req *PaymentRequest) (*PaymentResponse, error) {
	// In a real implementation, you would:
	// 1. Use the appropriate Xendit API endpoint based on payment method
	// 2. Create an invoice or payment request
	// 3. Return the payment response with URL, QR code, or VA number

	if x.logger != nil {
		x.logger.Info("creating Xendit transaction",
			zap.String("order_id", req.OrderID),
			zap.Int64("amount", req.Amount),
		)
	}

	// Placeholder implementation
	// In reality, you would call the Xendit API here
	return &PaymentResponse{
		TransactionID: fmt.Sprintf("xendit-%s", req.OrderID),
		OrderID:       req.OrderID,
		Status:        StatusPending,
		Amount:        req.Amount,
		PaymentURL:    fmt.Sprintf("https://checkout.xendit.co/v2/invoice/%s", req.OrderID),
		Metadata: map[string]interface{}{
			"gateway": "xendit",
		},
	}, nil
}

// GetTransaction retrieves a transaction by ID
func (x *XenditGateway) GetTransaction(ctx context.Context, transactionID string) (*PaymentResponse, error) {
	// In a real implementation, you would:
	// 1. Call the Xendit API to get invoice/payment status
	// 2. Map the status to our PaymentStatus
	// 3. Return the payment response

	if x.logger != nil {
		x.logger.Info("getting Xendit transaction",
			zap.String("transaction_id", transactionID),
		)
	}

	// Placeholder implementation
	return &PaymentResponse{
		TransactionID: transactionID,
		Status:        StatusPending,
	}, nil
}

// CancelTransaction cancels a transaction
func (x *XenditGateway) CancelTransaction(ctx context.Context, transactionID string) error {
	// In a real implementation, you would:
	// 1. Call the Xendit API to expire/cancel the invoice
	// 2. Handle the response

	if x.logger != nil {
		x.logger.Info("cancelling Xendit transaction",
			zap.String("transaction_id", transactionID),
		)
	}

	return nil
}

// VerifyNotification verifies a payment notification
func (x *XenditGateway) VerifyNotification(ctx context.Context, payload map[string]interface{}) (*PaymentNotification, error) {
	// In a real implementation, you would:
	// 1. Verify the webhook signature using x.config.WebhookToken
	// 2. Parse the webhook payload
	// 3. Return the notification details

	externalID, _ := payload["external_id"].(string)
	status, _ := payload["status"].(string)
	amount, _ := payload["amount"].(float64)

	paymentStatus := x.mapXenditStatus(status)

	var paidAt *time.Time
	if paymentStatus == StatusSuccess {
		if paidAtStr, ok := payload["paid_at"].(string); ok {
			t, err := time.Parse(time.RFC3339, paidAtStr)
			if err == nil {
				paidAt = &t
			}
		}
	}

	return &PaymentNotification{
		TransactionID: payload["id"].(string),
		OrderID:       externalID,
		Status:        paymentStatus,
		Amount:        int64(amount),
		PaidAt:        paidAt,
		Metadata:      payload,
	}, nil
}

// GetName returns the name of the payment gateway
func (x *XenditGateway) GetName() string {
	return "Xendit"
}

func (x *XenditGateway) mapXenditStatus(status string) PaymentStatus {
	switch status {
	case "PAID", "SETTLED":
		return StatusSuccess
	case "PENDING":
		return StatusPending
	case "EXPIRED":
		return StatusExpired
	case "FAILED":
		return StatusFailed
	default:
		return StatusPending
	}
}
