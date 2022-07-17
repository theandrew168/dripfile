package storage_test

import (
	"testing"

	"github.com/theandrew168/dripfile/internal/core"
	"github.com/theandrew168/dripfile/internal/test"
)

func mockTransfer(src, dst core.Location, schedule core.Schedule, project core.Project) core.Transfer {
	transfer := core.NewTransfer(
		test.RandomString(8),
		src,
		dst,
		schedule,
		project,
	)
	return transfer
}

func TestTransferCreate(t *testing.T) {
	store, closer := test.Storage(t)
	defer closer()

	project := mockProject()
	err := store.Project.Create(&project)
	if err != nil {
		t.Fatal(err)
	}

	schedule := mockSchedule(project)
	err = store.Schedule.Create(&schedule)
	if err != nil {
		t.Fatal(err)
	}

	src := mockLocation(project)
	err = store.Location.Create(&src)
	if err != nil {
		t.Fatal(err)
	}

	dst := mockLocation(project)
	err = store.Location.Create(&dst)
	if err != nil {
		t.Fatal(err)
	}

	transfer := mockTransfer(src, dst, schedule, project)
	err = store.Transfer.Create(&transfer)
	if err != nil {
		t.Fatal(err)
	}

	if transfer.ID == "" {
		t.Fatal("record ID should be non-empty after create")
	}
}
