package core

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
