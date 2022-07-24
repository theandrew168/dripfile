package task

type Kind string

type Status string

const (
	StatusNew     = "new"
	StatusRunning = "running"
	StatusError   = "error"
)

type Task struct {
	// readonly (from database, after creation)
	ID string

	Kind   Kind
	Info   string
	Status Status
	Error  string
}

func NewTask(kind Kind, info string) Task {
	t := Task{
		Kind:   kind,
		Info:   info,
		Status: StatusNew,
		Error:  "",
	}
	return t
}
