package core

// relationship between Transfers and Schedules
type Job struct {
	// readonly (from database, after creation)
	ID int64

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

type JobStorage interface {
	// baseline CRUD ops all deal with one record
	Create(job *Job) error
	Read(id int64) (Job, error)
	Update(job Job) error
	Delete(job Job) error
}
