package payment

import (
	"github.com/theandrew168/dripfile/internal/log"
	"github.com/theandrew168/dripfile/internal/random"
)

type logBilling struct {
	logger    log.Logger
	customers []string
}

func NewLogBilling(logger log.Logger) Billing {
	b := logBilling{
		logger: logger,
	}
	return &b
}

func (b *logBilling) CreateCustomer(email string) (string, error) {
	billingID := random.String(16)
	b.customers = append(b.customers, billingID)

	b.logger.Info("--- LogBilling Start ---\n")
	b.logger.Info("CreateCustomer:\n")
	b.logger.Info("BillingID: %s\n", billingID)
	b.logger.Info("--- LogBilling Finish ---\n")

	return billingID, nil
}

func (b *logBilling) CreateSetupIntent(billingID string) (string, error) {
	clientSecret := random.String(16)

	b.logger.Info("--- LogBilling Start ---\n")
	b.logger.Info("CreateSetupIntent:\n")
	b.logger.Info("BillingID: %s\n", billingID)
	b.logger.Info("ClientSecret: %s\n", clientSecret)
	b.logger.Info("--- LogBilling Finish ---\n")

	return clientSecret, nil
}
