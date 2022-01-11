package core

import (
	"github.com/theandrew168/dripfile/internal/core"
)

type Queue interface {
	Push(transfer core.Transfer) error
	Pop() (core.Transfer, error)
}
