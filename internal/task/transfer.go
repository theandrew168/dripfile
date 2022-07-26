package task

import (
	"context"
	"encoding/json"
	"errors"
	"time"

	"github.com/theandrew168/dripfile/internal/fileserver"
	"github.com/theandrew168/dripfile/internal/model"
	"github.com/theandrew168/dripfile/internal/postgresql"
)

const (
	KindTransferTry = "transfer:try"
)

type TransferTryInfo struct {
	TransferID string `json:"transfer_id"`
}

func NewTransferTryTask(transferID string) Task {
	info := TransferTryInfo{
		TransferID: transferID,
	}

	js, err := json.Marshal(info)
	if err != nil {
		panic(err)
	}

	return NewTask(KindTransferTry, string(js))
}

// TODO: write to tmp file and replace
func (w *Worker) HandleTransferTry(ctx context.Context, t Task) error {
	var info TransferTryInfo
	err := json.Unmarshal([]byte(t.Info), &info)
	if err != nil {
		return err
	}

	start := time.Now()

	// lookup transfer by ID
	transfer, err := w.store.Transfer.Read(info.TransferID)
	if err != nil {
		// transfer has since been deleted
		if errors.Is(err, postgresql.ErrNotExist) {
			w.logger.Info("transfer deleted", map[string]string{
				"transfer_id": transfer.ID,
			})
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

	// update history table
	finish := time.Now()
	history := model.NewHistory(
		totalBytes,
		"success",
		start,
		finish,
		transfer.ID,
		transfer.Project,
	)

	err = w.store.History.Create(&history)
	if err != nil {
		return err
	}

	return nil
}
