package task

type Queue interface {
	Push(task Task) error
	Pop() (Task, error)
	Update(task Task) error
}
