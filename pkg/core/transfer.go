package core

import (
	"github.com/theandrew168/dripfile/pkg/random"
)

type Transfer struct {
	// readonly (from database, after creation)
	ID string

	Pattern  string
	Src      Location
	Dst      Location
	Schedule Schedule
	Project  Project
}

func NewTransfer(pattern string, src, dst Location, schedule Schedule, project Project) Transfer {
	transfer := Transfer{
		Pattern:  pattern,
		Src:      src,
		Dst:      dst,
		Schedule: schedule,
		Project:  project,
	}
	return transfer
}

func NewTransferMock(src, dst Location, schedule Schedule, project Project) Transfer {
	transfer := NewTransfer(
		random.String(8),
		src,
		dst,
		schedule,
		project,
	)
	return transfer
}
