package transfer

type Transfer struct {
	// readonly (from database, after creation)
	ID string

	Pattern        string
	FromLocationID string
	ToLocationID   string
	ScheduleID     string
}

func NewTransfer(pattern, fromLocationID, toLocationID, scheduleID string) Transfer {
	transfer := Transfer{
		Pattern:        pattern,
		FromLocationID: fromLocationID,
		ToLocationID:   toLocationID,
		ScheduleID:     scheduleID,
	}
	return transfer
}
