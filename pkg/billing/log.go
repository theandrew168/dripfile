package billing

import (
	"github.com/theandrew168/dripfile/pkg/log"
	"github.com/theandrew168/dripfile/pkg/random"
)

type logGateway struct {
	logger    log.Logger
	customers []string
}

func NewLogGateway(logger log.Logger) PaymentGateway {
	g := logGateway{
		logger: logger,
	}
	return &g
}

func (g *logGateway) CreateCustomer(email string) (string, error) {
	billingID := random.String(16)
	g.customers = append(g.customers, billingID)

	g.logger.Info("--- LogGateway Start ---\n")
	g.logger.Info("CreateCustomer:\n")
	g.logger.Info("BillingID: %s\n", billingID)
	g.logger.Info("--- LogGateway Finish ---\n")

	return billingID, nil
}

func (g *logGateway) CreateCheckoutSession(billingID string) (string, error) {
	sessionURL := "/billing/success?session_id=" + random.String(16)

	g.logger.Info("--- LogGateway Start ---\n")
	g.logger.Info("CreateCheckoutSession:\n")
	g.logger.Info("BillingID: %s\n", billingID)
	g.logger.Info("SessionURL: %s\n", sessionURL)
	g.logger.Info("--- LogGateway Finish ---\n")

	return sessionURL, nil
}
