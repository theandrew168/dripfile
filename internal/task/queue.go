package task

type Queue interface {
	Publish(task Task) error
	Subscribe() (Task, error)
}
