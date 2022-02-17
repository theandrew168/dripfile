package pubsub

import (
	"github.com/theandrew168/dripfile/internal/core"
)

// aggregation of core queue interfaces
type Queue struct {
	Transfer TransferQueue
}

type TransferQueue interface {
	Publish(transfer core.Transfer) error
	Subscribe() (core.Transfer, error)
}
