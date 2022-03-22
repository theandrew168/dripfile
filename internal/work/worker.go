package work

import (
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/theandrew168/dripfile/internal/core"
	"github.com/theandrew168/dripfile/internal/database"
	"github.com/theandrew168/dripfile/internal/fileserver"
	"github.com/theandrew168/dripfile/internal/log"
	"github.com/theandrew168/dripfile/internal/mail"
	"github.com/theandrew168/dripfile/internal/secret"
	"github.com/theandrew168/dripfile/internal/task"
)

type TaskFunc func(t task.Task) error

type Worker struct {
	box     secret.Box
	queue   task.Queue
	storage database.Storage
	mailer  mail.Mailer
	logger  log.Logger
}

func NewWorker(box secret.Box, queue task.Queue, storage database.Storage, mailer mail.Mailer, logger log.Logger) *Worker {
	worker := Worker{
		box:     box,
		queue:   queue,
		storage: storage,
		mailer:  mailer,
		logger:  logger,
	}
	return &worker
}

// listen on queue, grab jobs, do the work, update as needed, success or error
func (w *Worker) Run() error {
	// check for new tasks periodically
	c := time.Tick(time.Second)
	for range c {
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
	}

	return nil
}

func (w *Worker) RunTask(t task.Task) {
	// determine which task needs to run
	var taskFunc TaskFunc
	switch t.Kind {
	case task.KindEmail:
		taskFunc = w.SendEmail
	case task.KindSession:
		taskFunc = w.DeleteExpiredSessions
	case task.KindTransfer:
		taskFunc = w.DoTransfer
	default:
		err := fmt.Errorf("unknown task: %s", t.Kind)
		w.logger.Error(err)
		return
	}

	w.logger.Info("task %s start\n", t.ID)

	// run and update the status
	err := taskFunc(t)
	if err != nil {
		w.logger.Error(err)

		t.Error = err.Error()
		t.Status = task.StatusError
	} else {
		t.Status = task.StatusSuccess
	}

	err = w.queue.Update(t)
	if err != nil {
		w.logger.Error(err)
	}

	w.logger.Info("task %s finish\n", t.ID)
}

func (w *Worker) SendEmail(t task.Task) error {
	var info task.SendEmailInfo
	err := json.Unmarshal([]byte(t.Info), &info)
	if err != nil {
		return err
	}

	err = w.mailer.SendEmail(
		info.FromName,
		info.FromEmail,
		info.ToName,
		info.ToEmail,
		info.Subject,
		info.Body,
	)
	if err != nil {
		return err
	}

	return nil
}

func (w *Worker) DeleteExpiredSessions(t task.Task) error {
	err := w.storage.Session.DeleteExpired()
	if err != nil {
		return err
	}

	return nil
}

func (w *Worker) DoTransfer(t task.Task) error {
	start := time.Now()

	var info task.DoTransferInfo
	err := json.Unmarshal([]byte(t.Info), &info)
	if err != nil {
		return err
	}

	// lookup transfer by ID
	transfer, err := w.storage.Transfer.Read(info.ID)
	if err != nil {
		return err
	}

	// decrypt src info
	src := transfer.Src
	srcBytes, err := w.box.Decrypt(src.Info)
	if err != nil {
		return err
	}

	// unmarshal src info json
	var srcInfo fileserver.S3Info
	err = json.Unmarshal(srcBytes, &srcInfo)
	if err != nil {
		return err
	}

	// create src fileserver
	srcConn, err := fileserver.NewS3(srcInfo)
	if err != nil {
		return err
	}

	// decrypt dst info
	dst := transfer.Dst
	dstBytes, err := w.box.Decrypt(dst.Info)
	if err != nil {
		return err
	}

	// unmarshal dst info json
	var dstInfo fileserver.S3Info
	err = json.Unmarshal(dstBytes, &dstInfo)
	if err != nil {
		return err
	}

	// create dst fileserver
	dstConn, err := fileserver.NewS3(dstInfo)
	if err != nil {
		return err
	}

	// search for matching files
	files, err := srcConn.Search(transfer.Pattern)
	if err != nil {
		return err
	}

	// transfer them all
	var total int64
	for _, file := range files {
		r, err := srcConn.Read(file)
		if err != nil {
			return err
		}

		err = dstConn.Write(file, r)
		if err != nil {
			return err
		}

		total += file.Size
	}

	// update history table (TODO: middleware?)
	finish := time.Now()
	history := core.NewHistory(
		total,
		"success",
		start,
		finish,
		transfer.ID,
		transfer.Project,
	)

	err = w.storage.History.Create(&history)
	if err != nil {
		return err
	}

	return nil
}
