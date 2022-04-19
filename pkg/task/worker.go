package task

import (
	"errors"
	"log"
	"time"

	"github.com/theandrew168/dripfile/pkg/core"
	"github.com/theandrew168/dripfile/pkg/postmark"
	"github.com/theandrew168/dripfile/pkg/secret"
	"github.com/theandrew168/dripfile/pkg/storage"
	"github.com/theandrew168/dripfile/pkg/stripe"
)

type TaskFunc func(task Task) error

type Worker struct {
	box      *secret.Box
	queue    *Queue
	storage  *storage.Storage
	stripe   stripe.Interface
	postmark postmark.Interface
	infoLog  *log.Logger
	errorLog *log.Logger
}

func NewWorker(
	box *secret.Box,
	queue *Queue,
	storage *storage.Storage,
	stripe stripe.Interface,
	postmark postmark.Interface,
	infoLog *log.Logger,
	errorLog *log.Logger,
) *Worker {
	worker := Worker{
		box:      box,
		queue:    queue,
		storage:  storage,
		stripe:   stripe,
		postmark: postmark,
		infoLog:  infoLog,
		errorLog: errorLog,
	}
	return &worker
}

// listen on queue, grab jobs, do the work, update as needed, success or error
func (w *Worker) Run() error {
	// check for new tasks periodically
	c := time.Tick(time.Second)
	for {
		// kick off all new tasks
		for {
			t, err := w.queue.Pop()
			if err != nil {
				// break loop if no new tasks remain
				if errors.Is(err, core.ErrNotExist) {
					break
				}
				return err
			}

			// run task in the background
			go w.RunTask(t)
		}

		<-c
	}

	return nil
}

func (w *Worker) RunTask(task Task) {
	// determine which task needs to run
	var taskFunc TaskFunc
	switch task.Kind {
	case KindPruneSessions:
		taskFunc = w.PruneSessions
	case KindSendEmail:
		taskFunc = w.SendEmail
	case KindTransfer:
		taskFunc = w.Transfer
	default:
		w.errorLog.Printf("unknown task: %s", task.Kind)
		return
	}

	w.infoLog.Printf("task %s start\n", task.ID)

	// run and update the status
	err := taskFunc(task)
	if err != nil {
		w.errorLog.Println(err)

		task.Error = err.Error()
		task.Status = StatusError
	} else {
		task.Status = StatusSuccess
	}

	err = w.queue.Update(task)
	if err != nil {
		w.errorLog.Println(err)
	}

	w.infoLog.Printf("task %s finish\n", task.ID)
}
