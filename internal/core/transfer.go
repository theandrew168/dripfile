package core

type Transfer struct {
	// readonly (from database, after creation)
	ID string

	Pattern string
	Src     Location
	Dst     Location
	Project Project
}

func NewTransfer(pattern string, src, dst Location, project Project) Transfer {
	transfer := Transfer{
		Pattern: pattern,
		Src:     src,
		Dst:     dst,
		Project: project,
	}
	return transfer
}
