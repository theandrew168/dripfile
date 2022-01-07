package core

import (
	"github.com/theandrew168/dripfile/internal/core"
)

type Queue interface {
	Push(job core.Job) error
	Pop() (core.Job, error)
}
