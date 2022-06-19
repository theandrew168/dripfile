package task

import (
	"context"
	"encoding/json"
	"errors"
	"math"
	"time"

	"github.com/hibiken/asynq"

	"github.com/theandrew168/dripfile/src/core"
	"github.com/theandrew168/dripfile/src/database"
	"github.com/theandrew168/dripfile/src/fileserver"
)

const (
	TypeTransferTry = "transfer:try"
)

type TransferTryPayload struct {
	TransferID string `json:"transfer_id"`
}

func NewTransferTryTask(transferID string) (*asynq.Task, error) {
	payload := TransferTryPayload{
		TransferID: transferID,
	}

	js, err := json.Marshal(payload)
	if err != nil {
		return nil, err
	}

	return asynq.NewTask(TypeTransferTry, js), nil
}

// TODO: write to tmp file and replace
func (w *Worker) HandleTransferTry(ctx context.Context, t *asynq.Task) error {
	var payload TransferTryPayload
	err := json.Unmarshal(t.Payload(), &payload)
	if err != nil {
		return err
	}

	start := time.Now()

	// lookup transfer by ID
	transfer, err := w.storage.Transfer.Read(payload.TransferID)
	if err != nil {
		// transfer has since been deleted
		if errors.Is(err, database.ErrNotExist) {
			return nil
		}
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
	// TODO: write all to temps, rename if success, else rollback
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
		totalBytes,
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
	customerID := transfer.Project.CustomerID
	subscriptionItemID := transfer.Project.SubscriptionItemID
	err = w.billing.CreateUsageRecord(customerID, subscriptionItemID, int64(mb))
	if err != nil {
		return err
	}

	return nil
}
