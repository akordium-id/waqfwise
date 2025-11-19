package payment

import (
	"context"
	"fmt"
	"time"

	"github.com/akordium-id/waqfwise/pkg/config"
	"github.com/midtrans/midtrans-go"
	"github.com/midtrans/midtrans-go/coreapi"
	"github.com/midtrans/midtrans-go/snap"
	"go.uber.org/zap"
)

// MidtransGateway implements PaymentGateway for Midtrans
type MidtransGateway struct {
	snapClient *snap.Client
	coreClient *coreapi.Client
	config     *config.MidtransConfig
	logger     *zap.Logger
}

// NewMidtransGateway creates a new Midtrans payment gateway
func NewMidtransGateway(cfg *config.MidtransConfig, logger *zap.Logger) *MidtransGateway {
	// Initialize Snap client
	snapClient := snap.Client{}
	snapClient.New(cfg.ServerKey, midtrans.Sandbox)
	if !cfg.IsSandbox() {
		snapClient.New(cfg.ServerKey, midtrans.Production)
	}

	// Initialize Core API client
	coreClient := coreapi.Client{}
	coreClient.New(cfg.ServerKey, midtrans.Sandbox)
	if !cfg.IsSandbox() {
		coreClient.New(cfg.ServerKey, midtrans.Production)
	}

	return &MidtransGateway{
		snapClient: &snapClient,
		coreClient: &coreClient,
		config:     cfg,
		logger:     logger,
	}
}

// CreateTransaction creates a new payment transaction
func (m *MidtransGateway) CreateTransaction(ctx context.Context, req *PaymentRequest) (*PaymentResponse, error) {
	// Build Snap request
	snapReq := &snap.Request{
		TransactionDetails: midtrans.TransactionDetails{
			OrderID:  req.OrderID,
			GrossAmt: req.Amount,
		},
		CustomerDetail: &midtrans.CustomerDetails{
			FName: req.CustomerName,
			Email: req.CustomerEmail,
			Phone: req.CustomerPhone,
		},
	}

	// Add items if provided
	if len(req.Items) > 0 {
		items := make([]midtrans.ItemDetails, len(req.Items))
		for i, item := range req.Items {
			items[i] = midtrans.ItemDetails{
				ID:    item.ID,
				Name:  item.Name,
				Price: item.Price,
				Qty:   int32(item.Quantity),
			}
		}
		snapReq.Items = &items
	}

	// Create transaction
	snapResp, err := m.snapClient.CreateTransaction(snapReq)
	if err != nil {
		if m.logger != nil {
			m.logger.Error("failed to create Midtrans transaction",
				zap.String("order_id", req.OrderID),
				zap.Error(err),
			)
		}
		return nil, fmt.Errorf("failed to create Midtrans transaction: %w", err)
	}

	if m.logger != nil {
		m.logger.Info("Midtrans transaction created",
			zap.String("order_id", req.OrderID),
			zap.String("token", snapResp.Token),
		)
	}

	return &PaymentResponse{
		TransactionID: snapResp.Token,
		OrderID:       req.OrderID,
		Status:        StatusPending,
		Amount:        req.Amount,
		PaymentURL:    snapResp.RedirectURL,
		Metadata: map[string]interface{}{
			"token": snapResp.Token,
		},
	}, nil
}

// GetTransaction retrieves a transaction by ID
func (m *MidtransGateway) GetTransaction(ctx context.Context, transactionID string) (*PaymentResponse, error) {
	// Get transaction status
	transactionStatusResp, err := m.coreClient.CheckTransaction(transactionID)
	if err != nil {
		return nil, fmt.Errorf("failed to get transaction status: %w", err)
	}

	status := m.mapMidtransStatus(transactionStatusResp.TransactionStatus)

	var paidAt *time.Time
	if status == StatusSuccess && transactionStatusResp.TransactionTime != "" {
		t, err := time.Parse("2006-01-02 15:04:05", transactionStatusResp.TransactionTime)
		if err == nil {
			paidAt = &t
		}
	}

	return &PaymentResponse{
		TransactionID: transactionStatusResp.TransactionID,
		OrderID:       transactionStatusResp.OrderID,
		Status:        status,
		Amount:        int64(transactionStatusResp.GrossAmount),
		PaidAt:        paidAt,
		Metadata: map[string]interface{}{
			"payment_type":    transactionStatusResp.PaymentType,
			"fraud_status":    transactionStatusResp.FraudStatus,
			"transaction_time": transactionStatusResp.TransactionTime,
		},
	}, nil
}

// CancelTransaction cancels a transaction
func (m *MidtransGateway) CancelTransaction(ctx context.Context, transactionID string) error {
	_, err := m.coreClient.CancelTransaction(transactionID)
	if err != nil {
		return fmt.Errorf("failed to cancel transaction: %w", err)
	}

	if m.logger != nil {
		m.logger.Info("Midtrans transaction cancelled",
			zap.String("transaction_id", transactionID),
		)
	}

	return nil
}

// VerifyNotification verifies a payment notification
func (m *MidtransGateway) VerifyNotification(ctx context.Context, payload map[string]interface{}) (*PaymentNotification, error) {
	// In a real implementation, you would:
	// 1. Verify the notification signature
	// 2. Parse the notification payload
	// 3. Return the notification details

	orderID, _ := payload["order_id"].(string)
	transactionStatus, _ := payload["transaction_status"].(string)
	grossAmount, _ := payload["gross_amount"].(string)

	status := m.mapMidtransStatus(transactionStatus)

	var amount int64
	fmt.Sscanf(grossAmount, "%d", &amount)

	var paidAt *time.Time
	if status == StatusSuccess {
		now := time.Now()
		paidAt = &now
	}

	return &PaymentNotification{
		TransactionID: payload["transaction_id"].(string),
		OrderID:       orderID,
		Status:        status,
		Amount:        amount,
		PaidAt:        paidAt,
		Metadata:      payload,
	}, nil
}

// GetName returns the name of the payment gateway
func (m *MidtransGateway) GetName() string {
	return "Midtrans"
}

func (m *MidtransGateway) mapMidtransStatus(status string) PaymentStatus {
	switch status {
	case "capture", "settlement":
		return StatusSuccess
	case "pending":
		return StatusPending
	case "deny", "cancel":
		return StatusCancelled
	case "expire":
		return StatusExpired
	default:
		return StatusFailed
	}
}
