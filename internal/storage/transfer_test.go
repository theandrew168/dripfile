package storage_test

import (
	"testing"

	"github.com/theandrew168/dripfile/internal/model"
	"github.com/theandrew168/dripfile/internal/test"
)

func mockTransfer(src, dst model.Location, schedule model.Schedule, project model.Project) model.Transfer {
	transfer := model.NewTransfer(
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

	err = store.Project.Delete(project)
	if err != nil {
		t.Fatal(err)
	}
}
