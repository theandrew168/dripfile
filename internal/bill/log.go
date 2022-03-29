package bill

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
	id := random.String(16)
	b.customers = append(b.customers, id)

	b.logger.Info("--- LogBilling Start ---\n")
	b.logger.Info("CreateCustomer:\n")
	b.logger.Info("ID: %s\n", id)
	b.logger.Info("--- LogBilling Finish ---\n")

	return id, nil
}
