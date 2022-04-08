package task

import (
	"encoding/json"
	"errors"
	"log"
	"math"
	"time"

	"github.com/theandrew168/dripfile/pkg/core"
	"github.com/theandrew168/dripfile/pkg/fileserver"
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

func (w *Worker) RunTask(task Task) {
	// determine which task needs to run
	var taskFunc TaskFunc
	switch task.Kind {
	case KindEmail:
		taskFunc = w.SendEmail
	case KindSession:
		taskFunc = w.DeleteExpiredSessions
	case KindTransfer:
		taskFunc = w.DoTransfer
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

func (w *Worker) SendEmail(task Task) error {
	var info SendEmailInfo
	err := json.Unmarshal([]byte(task.Info), &info)
	if err != nil {
		return err
	}

	err = w.postmark.SendEmail(
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

func (w *Worker) DeleteExpiredSessions(task Task) error {
	err := w.storage.Session.DeleteExpired()
	if err != nil {
		return err
	}

	return nil
}

func (w *Worker) DoTransfer(task Task) error {
	start := time.Now()

	var info DoTransferInfo
	err := json.Unmarshal([]byte(task.Info), &info)
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
	var totalBytes int64
	for _, file := range files {
		r, err := srcConn.Read(file)
		if err != nil {
			return err
		}

		err = dstConn.Write(file, r)
		if err != nil {
			return err
		}

		totalBytes += file.Size
	}

	// convert total bytes to megabytes
	mb := math.Ceil(float64(totalBytes) / (1000 * 1000))

	// update history table
	finish := time.Now()
	history := core.NewHistory(
		int64(mb),
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

	// create usage record
	subscriptionItemID := transfer.Project.SubscriptionItemID
	err = w.stripe.CreateUsageRecord(subscriptionItemID, int64(mb))
	if err != nil {
		return err
	}

	return nil
}
