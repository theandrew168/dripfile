package storage_test

import (
	"errors"
	"testing"

	"github.com/theandrew168/dripfile/pkg/core"
	"github.com/theandrew168/dripfile/pkg/test"
)

func TestCreateProject(t *testing.T) {
	storage, closer := test.Storage(t)
	defer closer()

	project := core.NewMockProject()
	err := storage.Project.Create(&project)
	if err != nil {
		t.Fatal(err)
	}

	if project.ID == "" {
		t.Fatal("record ID should be non-empty after create")
	}
}

func TestCreateProjectDuplicate(t *testing.T) {
	storage, closer := test.Storage(t)
	defer closer()

	project := core.NewMockProject()
	err := storage.Project.Create(&project)
	if err != nil {
		t.Fatal(err)
	}

	// attempt to create the same project again
	err = storage.Project.Create(&project)
	if !errors.Is(err, core.ErrExist) {
		t.Fatal("duplicate record should return an error")
	}
}
