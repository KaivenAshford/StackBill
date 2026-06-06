package service

import (
	"bytes"
	"encoding/csv"
	"fmt"

	"github.com/kingqaquuu/stackbill/internal/repository"
)

type ExportService struct {
	subscriptionRepo *repository.SubscriptionRepository
	assetRepo        *repository.AssetRepository
}

func NewExportService(
	subscriptionRepo *repository.SubscriptionRepository,
	assetRepo *repository.AssetRepository,
) *ExportService {
	return &ExportService{
		subscriptionRepo: subscriptionRepo,
		assetRepo:        assetRepo,
	}
}

func (s *ExportService) ExportSubscriptionsCSV(userID uint) ([]byte, error) {
	subs, err := s.subscriptionRepo.GetActiveByUserID(userID)
	if err != nil {
		return nil, err
	}

	var buf bytes.Buffer
	w := csv.NewWriter(&buf)

	headers := []string{"name", "description", "category_id", "amount", "currency", "billing_cycle", "billing_interval", "start_date", "next_payment_date", "payment_method", "auto_renew", "status", "website_url", "remark"}
	if err := w.Write(headers); err != nil {
		return nil, err
	}

	for _, sub := range subs {
		autoRenew := "false"
		if sub.AutoRenew {
			autoRenew = "true"
		}
		startDate := ""
		if sub.StartDate != nil {
			startDate = sub.StartDate.Format("2006-01-02")
		}
		nextPayment := ""
		if sub.NextPaymentDate != nil {
			nextPayment = sub.NextPaymentDate.Format("2006-01-02")
		}
		row := []string{
			sub.Name, sub.Description, fmt.Sprintf("%d", sub.CategoryID),
			fmt.Sprintf("%.2f", sub.Amount), sub.Currency, sub.BillingCycle,
			fmt.Sprintf("%d", sub.BillingInterval), startDate, nextPayment,
			sub.PaymentMethod, autoRenew, sub.Status, sub.WebsiteURL, sub.Remark,
		}
		if err := w.Write(row); err != nil {
			return nil, err
		}
	}

	w.Flush()
	return buf.Bytes(), w.Error()
}

func (s *ExportService) ExportAssetsCSV(userID uint) ([]byte, error) {
	assets, err := s.assetRepo.GetByUserID(userID)
	if err != nil {
		return nil, err
	}

	var buf bytes.Buffer
	w := csv.NewWriter(&buf)

	headers := []string{"name", "asset_type", "provider", "identifier", "url", "expire_date", "cost_amount", "cost_currency", "billing_cycle", "status", "description", "remark", "subscription_id"}
	if err := w.Write(headers); err != nil {
		return nil, err
	}

	for _, asset := range assets {
		expireDate := ""
		if asset.ExpireDate != nil {
			expireDate = asset.ExpireDate.Format("2006-01-02")
		}
		row := []string{
			asset.Name, asset.AssetType, asset.Provider, asset.Identifier,
			asset.URL, expireDate, fmt.Sprintf("%.2f", asset.CostAmount),
			asset.CostCurrency, asset.BillingCycle, asset.Status,
			asset.Description, asset.Remark, fmt.Sprintf("%d", asset.SubscriptionID),
		}
		if err := w.Write(row); err != nil {
			return nil, err
		}
	}

	w.Flush()
	return buf.Bytes(), w.Error()
}
