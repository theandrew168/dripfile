package task

type Task struct {
	Kind string
	Info string
}

func NewTask(kind, info string) Task {
	t := Task{
		Kind: kind,
		Info: info,
	}
	return t
}
