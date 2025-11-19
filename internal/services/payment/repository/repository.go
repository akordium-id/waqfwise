package repository

import (
	"context"
	"database/sql"
	"time"

	"github.com/akordium-id/waqfwise/internal/shared/domain"
	"github.com/akordium-id/waqfwise/internal/shared/errors"
)

// Repository defines payment repository interface
type Repository interface {
	CreateDonation(ctx context.Context, donation *domain.Donation) error
	FindDonationByID(ctx context.Context, id int64) (*domain.Donation, error)
	FindDonationByTransactionID(ctx context.Context, txID string) (*domain.Donation, error)
	UpdateDonationStatus(ctx context.Context, id int64, status domain.PaymentStatus) error
	CreatePaymentLog(ctx context.Context, log *domain.PaymentLog) error
	CreateLedgerEntry(ctx context.Context, entry *domain.Ledger) error
	GetCampaignBalance(ctx context.Context, campaignID int64) (int64, error)
	CreateFraudCheck(ctx context.Context, check *domain.FraudCheck) error
	GetDonationsByUser(ctx context.Context, userID int64, limit, offset int) ([]*domain.Donation, int64, error)
	GetDonationsByCampaign(ctx context.Context, campaignID int64, limit, offset int) ([]*domain.Donation, int64, error)
}

type repository struct {
	db *sql.DB
}

// New creates a new payment repository
func New(db *sql.DB) Repository {
	return &repository{db: db}
}

// CreateDonation creates a new donation
func (r *repository) CreateDonation(ctx context.Context, donation *domain.Donation) error {
	query := `
		INSERT INTO donations (campaign_id, user_id, amount, status, payment_method, payment_gateway,
		                       transaction_id, is_anonymous, donor_name, donor_email, message,
		                       is_recurring, recurring_period, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15)
		RETURNING id
	`

	now := time.Now()
	err := r.db.QueryRowContext(
		ctx, query,
		donation.CampaignID,
		donation.UserID,
		donation.Amount,
		donation.Status,
		donation.PaymentMethod,
		donation.PaymentGateway,
		donation.TransactionID,
		donation.IsAnonymous,
		donation.DonorName,
		donation.DonorEmail,
		donation.Message,
		donation.IsRecurring,
		donation.RecurringPeriod,
		now,
		now,
	).Scan(&donation.ID)

	if err != nil {
		return errors.Wrap(err, errors.ErrCodeInternal, "Failed to create donation", 500)
	}

	donation.CreatedAt = now
	donation.UpdatedAt = now
	return nil
}

// FindDonationByID finds donation by ID
func (r *repository) FindDonationByID(ctx context.Context, id int64) (*domain.Donation, error) {
	query := `
		SELECT id, campaign_id, user_id, amount, status, payment_method, payment_gateway,
		       transaction_id, gateway_ref, is_anonymous, donor_name, donor_email, message,
		       is_recurring, recurring_period, receipt_url, paid_at, created_at, updated_at
		FROM donations
		WHERE id = $1
	`

	donation := &domain.Donation{}
	var paidAt sql.NullTime
	var gatewayRef, donorName, donorEmail, message, recurringPeriod, receiptURL sql.NullString

	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&donation.ID,
		&donation.CampaignID,
		&donation.UserID,
		&donation.Amount,
		&donation.Status,
		&donation.PaymentMethod,
		&donation.PaymentGateway,
		&donation.TransactionID,
		&gatewayRef,
		&donation.IsAnonymous,
		&donorName,
		&donorEmail,
		&message,
		&donation.IsRecurring,
		&recurringPeriod,
		&receiptURL,
		&paidAt,
		&donation.CreatedAt,
		&donation.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, errors.New(errors.ErrCodeNotFound, "Donation not found", 404)
	}
	if err != nil {
		return nil, errors.Wrap(err, errors.ErrCodeInternal, "Failed to find donation", 500)
	}

	if gatewayRef.Valid {
		donation.GatewayRef = gatewayRef.String
	}
	if donorName.Valid {
		donation.DonorName = donorName.String
	}
	if donorEmail.Valid {
		donation.DonorEmail = donorEmail.String
	}
	if message.Valid {
		donation.Message = message.String
	}
	if recurringPeriod.Valid {
		donation.RecurringPeriod = recurringPeriod.String
	}
	if receiptURL.Valid {
		donation.ReceiptURL = receiptURL.String
	}
	if paidAt.Valid {
		donation.PaidAt = &paidAt.Time
	}

	return donation, nil
}

// FindDonationByTransactionID finds donation by transaction ID
func (r *repository) FindDonationByTransactionID(ctx context.Context, txID string) (*domain.Donation, error) {
	query := `
		SELECT id, campaign_id, user_id, amount, status, payment_method, payment_gateway,
		       transaction_id, gateway_ref, is_anonymous, donor_name, donor_email, message,
		       is_recurring, recurring_period, receipt_url, paid_at, created_at, updated_at
		FROM donations
		WHERE transaction_id = $1
	`

	donation := &domain.Donation{}
	var paidAt sql.NullTime
	var gatewayRef, donorName, donorEmail, message, recurringPeriod, receiptURL sql.NullString

	err := r.db.QueryRowContext(ctx, query, txID).Scan(
		&donation.ID,
		&donation.CampaignID,
		&donation.UserID,
		&donation.Amount,
		&donation.Status,
		&donation.PaymentMethod,
		&donation.PaymentGateway,
		&donation.TransactionID,
		&gatewayRef,
		&donation.IsAnonymous,
		&donorName,
		&donorEmail,
		&message,
		&donation.IsRecurring,
		&recurringPeriod,
		&receiptURL,
		&paidAt,
		&donation.CreatedAt,
		&donation.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, errors.New(errors.ErrCodeNotFound, "Donation not found", 404)
	}
	if err != nil {
		return nil, errors.Wrap(err, errors.ErrCodeInternal, "Failed to find donation", 500)
	}

	// Handle nullable fields
	if gatewayRef.Valid {
		donation.GatewayRef = gatewayRef.String
	}
	if donorName.Valid {
		donation.DonorName = donorName.String
	}
	if donorEmail.Valid {
		donation.DonorEmail = donorEmail.String
	}
	if message.Valid {
		donation.Message = message.String
	}
	if recurringPeriod.Valid {
		donation.RecurringPeriod = recurringPeriod.String
	}
	if receiptURL.Valid {
		donation.ReceiptURL = receiptURL.String
	}
	if paidAt.Valid {
		donation.PaidAt = &paidAt.Time
	}

	return donation, nil
}

// UpdateDonationStatus updates donation status
func (r *repository) UpdateDonationStatus(ctx context.Context, id int64, status domain.PaymentStatus) error {
	query := `UPDATE donations SET status = $1, updated_at = $2 WHERE id = $3`

	result, err := r.db.ExecContext(ctx, query, status, time.Now(), id)
	if err != nil {
		return errors.Wrap(err, errors.ErrCodeInternal, "Failed to update donation status", 500)
	}

	rows, _ := result.RowsAffected()
	if rows == 0 {
		return errors.New(errors.ErrCodeNotFound, "Donation not found", 404)
	}

	return nil
}

// CreatePaymentLog creates payment log
func (r *repository) CreatePaymentLog(ctx context.Context, log *domain.PaymentLog) error {
	query := `
		INSERT INTO payment_logs (donation_id, status, gateway, request_data, response_data,
		                          error_message, ip_address, user_agent, created_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
		RETURNING id
	`

	err := r.db.QueryRowContext(
		ctx, query,
		log.DonationID,
		log.Status,
		log.Gateway,
		log.RequestData,
		log.ResponseData,
		log.ErrorMessage,
		log.IPAddress,
		log.UserAgent,
		time.Now(),
	).Scan(&log.ID)

	if err != nil {
		return errors.Wrap(err, errors.ErrCodeInternal, "Failed to create payment log", 500)
	}

	return nil
}

// CreateLedgerEntry creates ledger entry
func (r *repository) CreateLedgerEntry(ctx context.Context, entry *domain.Ledger) error {
	query := `
		INSERT INTO ledgers (donation_id, account_type, account_name, amount, balance_before,
		                     balance_after, description, created_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
		RETURNING id
	`

	err := r.db.QueryRowContext(
		ctx, query,
		entry.DonationID,
		entry.AccountType,
		entry.AccountName,
		entry.Amount,
		entry.BalanceBefore,
		entry.BalanceAfter,
		entry.Description,
		time.Now(),
	).Scan(&entry.ID)

	if err != nil {
		return errors.Wrap(err, errors.ErrCodeInternal, "Failed to create ledger entry", 500)
	}

	return nil
}

// GetCampaignBalance gets campaign balance from ledger
func (r *repository) GetCampaignBalance(ctx context.Context, campaignID int64) (int64, error) {
	query := `
		SELECT COALESCE(SUM(CASE WHEN account_type = 'credit' THEN amount ELSE -amount END), 0)
		FROM ledgers l
		JOIN donations d ON l.donation_id = d.id
		WHERE d.campaign_id = $1
	`

	var balance int64
	err := r.db.QueryRowContext(ctx, query, campaignID).Scan(&balance)
	if err != nil {
		return 0, errors.Wrap(err, errors.ErrCodeInternal, "Failed to get campaign balance", 500)
	}

	return balance, nil
}

// CreateFraudCheck creates fraud check record
func (r *repository) CreateFraudCheck(ctx context.Context, check *domain.FraudCheck) error {
	query := `
		INSERT INTO fraud_checks (donation_id, risk_score, risk_level, flags, is_blocked,
		                          reason, ip_address, device_id, created_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
		RETURNING id
	`

	err := r.db.QueryRowContext(
		ctx, query,
		check.DonationID,
		check.RiskScore,
		check.RiskLevel,
		check.Flags,
		check.IsBlocked,
		check.Reason,
		check.IPAddress,
		check.DeviceID,
		time.Now(),
	).Scan(&check.ID)

	if err != nil {
		return errors.Wrap(err, errors.ErrCodeInternal, "Failed to create fraud check", 500)
	}

	return nil
}

// GetDonationsByUser gets donations by user
func (r *repository) GetDonationsByUser(ctx context.Context, userID int64, limit, offset int) ([]*domain.Donation, int64, error) {
	// Get total count
	var total int64
	countQuery := `SELECT COUNT(*) FROM donations WHERE user_id = $1`
	if err := r.db.QueryRowContext(ctx, countQuery, userID).Scan(&total); err != nil {
		return nil, 0, errors.Wrap(err, errors.ErrCodeInternal, "Failed to count donations", 500)
	}

	// Get donations
	query := `
		SELECT id, campaign_id, user_id, amount, status, payment_method, payment_gateway,
		       transaction_id, is_anonymous, message, created_at
		FROM donations
		WHERE user_id = $1
		ORDER BY created_at DESC
		LIMIT $2 OFFSET $3
	`

	rows, err := r.db.QueryContext(ctx, query, userID, limit, offset)
	if err != nil {
		return nil, 0, errors.Wrap(err, errors.ErrCodeInternal, "Failed to get donations", 500)
	}
	defer rows.Close()

	donations := make([]*domain.Donation, 0)
	for rows.Next() {
		d := &domain.Donation{}
		var message sql.NullString

		if err := rows.Scan(
			&d.ID, &d.CampaignID, &d.UserID, &d.Amount, &d.Status,
			&d.PaymentMethod, &d.PaymentGateway, &d.TransactionID,
			&d.IsAnonymous, &message, &d.CreatedAt,
		); err != nil {
			return nil, 0, errors.Wrap(err, errors.ErrCodeInternal, "Failed to scan donation", 500)
		}

		if message.Valid {
			d.Message = message.String
		}

		donations = append(donations, d)
	}

	return donations, total, nil
}

// GetDonationsByCampaign gets donations by campaign
func (r *repository) GetDonationsByCampaign(ctx context.Context, campaignID int64, limit, offset int) ([]*domain.Donation, int64, error) {
	// Get total count
	var total int64
	countQuery := `SELECT COUNT(*) FROM donations WHERE campaign_id = $1 AND status = $2`
	if err := r.db.QueryRowContext(ctx, countQuery, campaignID, domain.PaymentStatusSuccess).Scan(&total); err != nil {
		return nil, 0, errors.Wrap(err, errors.ErrCodeInternal, "Failed to count donations", 500)
	}

	// Get donations
	query := `
		SELECT id, campaign_id, user_id, amount, status, payment_method, payment_gateway,
		       transaction_id, is_anonymous, donor_name, message, created_at
		FROM donations
		WHERE campaign_id = $1 AND status = $2
		ORDER BY created_at DESC
		LIMIT $3 OFFSET $4
	`

	rows, err := r.db.QueryContext(ctx, query, campaignID, domain.PaymentStatusSuccess, limit, offset)
	if err != nil {
		return nil, 0, errors.Wrap(err, errors.ErrCodeInternal, "Failed to get donations", 500)
	}
	defer rows.Close()

	donations := make([]*domain.Donation, 0)
	for rows.Next() {
		d := &domain.Donation{}
		var donorName, message sql.NullString

		if err := rows.Scan(
			&d.ID, &d.CampaignID, &d.UserID, &d.Amount, &d.Status,
			&d.PaymentMethod, &d.PaymentGateway, &d.TransactionID,
			&d.IsAnonymous, &donorName, &message, &d.CreatedAt,
		); err != nil {
			return nil, 0, errors.Wrap(err, errors.ErrCodeInternal, "Failed to scan donation", 500)
		}

		if donorName.Valid {
			d.DonorName = donorName.String
		}
		if message.Valid {
			d.Message = message.String
		}

		donations = append(donations, d)
	}

	return donations, total, nil
}
