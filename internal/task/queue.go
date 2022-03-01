package task

import (
	"context"
)

type Queue interface {
	Push(task Task) error
	Pop() (Task, error)
	Listen(ctx context.Context) error
}
