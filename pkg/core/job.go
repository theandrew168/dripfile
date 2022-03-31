package core

// relationship between Transfers and Schedules
type Job struct {
	// readonly (from database, after creation)
	ID string

	Transfer Transfer
	Schedule Schedule
}

func NewJob(transfer Transfer, schedule Schedule) Job {
	job := Job{
		Transfer: transfer,
		Schedule: schedule,
	}
	return job
}
