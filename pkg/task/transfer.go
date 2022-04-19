package task

import (
	"encoding/json"
	"errors"
	"math"
	"time"

	"github.com/theandrew168/dripfile/pkg/core"
	"github.com/theandrew168/dripfile/pkg/fileserver"
)

const KindTransfer = "transfer"

type TransferInfo struct {
	ID string `json:"id"`
}

func Transfer(id string) (Task, error) {
	info := TransferInfo{
		ID: id,
	}

	b, err := json.Marshal(info)
	if err != nil {
		return Task{}, err
	}

	return New(KindTransfer, string(b)), nil
}

// TODO: write to tmp file and replace
func (w *Worker) Transfer(task Task) error {
	start := time.Now()

	var info TransferInfo
	err := json.Unmarshal([]byte(task.Info), &info)
	if err != nil {
		return err
	}

	// lookup transfer by ID
	transfer, err := w.storage.Transfer.Read(info.ID)
	if err != nil {
		// transfer has since been deleted
		if errors.Is(err, core.ErrNotExist) {
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
	subscriptionItemID := transfer.Project.SubscriptionItemID
	err = w.stripe.CreateUsageRecord(subscriptionItemID, int64(mb))
	if err != nil {
		return err
	}

	return nil
}
