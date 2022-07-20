package task

import (
	"encoding/json"
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
