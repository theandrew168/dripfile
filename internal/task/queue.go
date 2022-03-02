package task

type Queue interface {
	Push(task Task) error
	Pop() (Task, error)
}
