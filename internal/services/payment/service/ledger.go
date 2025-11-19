package service

import (
	"context"
	"fmt"

	"github.com/akordium-id/waqfwise/internal/services/payment/repository"
	"github.com/akordium-id/waqfwise/internal/shared/domain"
)

// LedgerManager handles double-entry bookkeeping
type LedgerManager struct {
	repo repository.Repository
}

// NewLedgerManager creates a new ledger manager
func NewLedgerManager(repo repository.Repository) *LedgerManager {
	return &LedgerManager{repo: repo}
}

// RecordDonation records donation in ledger using double-entry bookkeeping
func (l *LedgerManager) RecordDonation(ctx context.Context, donation *domain.Donation) error {
	// Calculate amounts
	totalAmount := donation.Amount
	gatewayFee := l.calculateGatewayFee(donation.PaymentGateway, totalAmount)
	netAmount := totalAmount - gatewayFee

	// Get current campaign balance
	currentBalance, err := l.repo.GetCampaignBalance(ctx, donation.CampaignID)
	if err != nil {
		return err
	}

	// Entry 1: Debit - Cash/Bank Account (money coming in)
	debitEntry := &domain.Ledger{
		DonationID:    donation.ID,
		AccountType:   "debit",
		AccountName:   "cash_account",
		Amount:        totalAmount,
		BalanceBefore: currentBalance,
		BalanceAfter:  currentBalance + totalAmount,
		Description:   fmt.Sprintf("Donation received from transaction %s", donation.TransactionID),
	}

	if err := l.repo.CreateLedgerEntry(ctx, debitEntry); err != nil {
		return err
	}

	// Entry 2: Credit - Campaign Fund (liability to campaign)
	creditEntry := &domain.Ledger{
		DonationID:    donation.ID,
		AccountType:   "credit",
		AccountName:   "campaign_fund",
		Amount:        netAmount,
		BalanceBefore: currentBalance,
		BalanceAfter:  currentBalance + netAmount,
		Description:   fmt.Sprintf("Campaign fund for campaign ID %d", donation.CampaignID),
	}

	if err := l.repo.CreateLedgerEntry(ctx, creditEntry); err != nil {
		return err
	}

	// Entry 3: Credit - Gateway Fee Expense
	if gatewayFee > 0 {
		feeEntry := &domain.Ledger{
			DonationID:    donation.ID,
			AccountType:   "credit",
			AccountName:   "gateway_fee_expense",
			Amount:        gatewayFee,
			BalanceBefore: currentBalance,
			BalanceAfter:  currentBalance,
			Description:   fmt.Sprintf("Gateway fee for %s", donation.PaymentGateway),
		}

		if err := l.repo.CreateLedgerEntry(ctx, feeEntry); err != nil {
			return err
		}
	}

	return nil
}

// calculateGatewayFee calculates payment gateway fee
func (l *LedgerManager) calculateGatewayFee(gateway domain.PaymentGateway, amount int64) int64 {
	switch gateway {
	case domain.PaymentGatewayMidtrans:
		// Midtrans: 2% + Rp 2000
		return (amount * 2 / 100) + 2000
	case domain.PaymentGatewayXendit:
		// Xendit: 2.9% + Rp 2000
		return (amount * 29 / 1000) + 2000
	default:
		return 0
	}
}
