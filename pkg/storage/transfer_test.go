package storage_test

import (
	"testing"

	"github.com/theandrew168/dripfile/pkg/core"
	"github.com/theandrew168/dripfile/pkg/test"
)

func TestTransferCreate(t *testing.T) {
	storage, closer := test.Storage(t)
	defer closer()

	project := core.NewProjectMock()
	err := storage.Project.Create(&project)
	if err != nil {
		t.Fatal(err)
	}

	schedule := core.NewScheduleMock(project)
	err = storage.Schedule.Create(&schedule)
	if err != nil {
		t.Fatal(err)
	}

	src := core.NewLocationMock(project)
	err = storage.Location.Create(&src)
	if err != nil {
		t.Fatal(err)
	}

	dst := core.NewLocationMock(project)
	err = storage.Location.Create(&dst)
	if err != nil {
		t.Fatal(err)
	}

	transfer := core.NewTransferMock(src, dst, schedule, project)
	err = storage.Transfer.Create(&transfer)
	if err != nil {
		t.Fatal(err)
	}

	if transfer.ID == "" {
		t.Fatal("record ID should be non-empty after create")
	}
}
