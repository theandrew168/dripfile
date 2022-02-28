package task

const (
	KindEmail    = "email"
	KindSession  = "session"
	KindTransfer = "transfer"
)

const (
	StatusNew     = "new"
	StatusRunning = "running"
	StatusSuccess = "success"
	StatusError   = "error"
)

type Task struct {
	// readonly (from database, after creation)
	ID string

	Kind   string
	Info   string
	Status string
}

func New(kind, info string) Task {
	task := Task{
		Kind:   kind,
		Info:   info,
		Status: StatusNew,
	}
	return task
}
