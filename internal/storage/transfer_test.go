package storage_test

import (
	"github.com/theandrew168/dripfile/internal/model"
	"github.com/theandrew168/dripfile/internal/test"
)

func mockTransfer(src, dst model.Location, schedule model.Schedule) model.Transfer {
	transfer := model.NewTransfer(
		test.RandomString(8),
		src,
		dst,
		schedule,
	)
	return transfer
}

//func TestTransferCreate(t *testing.T) {
//	store, closer := test.Storage(t)
//	defer closer()
//
//	schedule := mockSchedule()
//	err := store.Schedule.Create(&schedule)
//	if err != nil {
//		t.Fatal(err)
//	}
//
//	src := mockLocation()
//	err = store.Location.Create(&src)
//	if err != nil {
//		t.Fatal(err)
//	}
//
//	dst := mockLocation()
//	err = store.Location.Create(&dst)
//	if err != nil {
//		t.Fatal(err)
//	}
//
//	transfer := mockTransfer(src, dst, schedule)
//	err = store.Transfer.Create(&transfer)
//	if err != nil {
//		t.Fatal(err)
//	}
//
//	if transfer.ID == "" {
//		t.Fatal("record ID should be non-empty after create")
//	}
//}
