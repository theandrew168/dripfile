package model

type Transfer struct {
	// readonly (from database, after creation)
	ID string

	Pattern  string
	Src      Location
	Dst      Location
	Schedule Schedule
}

func NewTransfer(pattern string, src, dst Location, schedule Schedule) Transfer {
	transfer := Transfer{
		Pattern:  pattern,
		Src:      src,
		Dst:      dst,
		Schedule: schedule,
	}
	return transfer
}
