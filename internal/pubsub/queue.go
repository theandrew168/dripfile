package pubsub

import (
	"github.com/theandrew168/dripfile/internal/core"
)

// aggregation of core queue interfaces
type Queue struct {
	Transfer TransferQueue
}

type TransferQueue interface {
	Push(transfer core.Transfer) error
	Pop() (core.Transfer, error)
}
