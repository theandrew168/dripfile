package storage_test

import (
	"errors"
	"testing"

	"github.com/theandrew168/dripfile/pkg/core"
	"github.com/theandrew168/dripfile/pkg/random"
	"github.com/theandrew168/dripfile/pkg/test"
)

func TestProjectCreate(t *testing.T) {
	storage, closer := test.Storage(t)
	defer closer()

	project := core.NewProjectMock()
	err := storage.Project.Create(&project)
	if err != nil {
		t.Fatal(err)
	}

	if project.ID == "" {
		t.Fatal("record ID should be non-empty after create")
	}
}

func TestProjectCreateUnique(t *testing.T) {
	storage, closer := test.Storage(t)
	defer closer()

	project := core.NewProjectMock()
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

func TestProjectRead(t *testing.T) {
	storage, closer := test.Storage(t)
	defer closer()

	project := core.NewProjectMock()
	err := storage.Project.Create(&project)
	if err != nil {
		t.Fatal(err)
	}

	got, err := storage.Project.Read(project.ID)
	if err != nil {
		t.Fatal(err)
	}
	
	if got.ID != project.ID {
		t.Fatal("record ID should be match after read")
	}
}

func TestProjectUpdate(t *testing.T) {
	storage, closer := test.Storage(t)
	defer closer()

	project := core.NewProjectMock()
	err := storage.Project.Create(&project)
	if err != nil {
		t.Fatal(err)
	}

	// update fields to new random values
	project.CustomerID = random.String(8)
	project.SubscriptionItemID = random.String(8)

	err = storage.Project.Update(project)
	if err != nil {
		t.Fatal(err)
	}

	got, err := storage.Project.Read(project.ID)
	if err != nil {
		t.Fatal(err)
	}

	if got.CustomerID != project.CustomerID {
		t.Fatal("CustomerID should be match after update")
	}
	if got.SubscriptionItemID != project.SubscriptionItemID {
		t.Fatal("SubscriptionItemID should be match after update")
	}
}

func TestProjectDelete(t *testing.T) {
	storage, closer := test.Storage(t)
	defer closer()

	project := core.NewProjectMock()
	err := storage.Project.Create(&project)
	if err != nil {
		t.Fatal(err)
	}

	err = storage.Project.Delete(project)
	if err != nil {
		t.Fatal(err)
	}

	// verify that ID isn't present anymore
	_, err = storage.Project.Read(project.ID)
	if !errors.Is(err, core.ErrNotExist) {
		t.Fatal("record ID should be gone after delete")
	}
}

func TestProjectReadAll(t *testing.T) {
	storage, closer := test.Storage(t)
	defer closer()

	project := core.NewProjectMock()
	err := storage.Project.Create(&project)
	if err != nil {
		t.Fatal(err)
	}

	projects, err := storage.Project.ReadAll()
	if err != nil {
		t.Fatal(err)
	}

	if len(projects) < 1 {
		t.Fatal("at least one record should be present")
	}

	// verify that ID is present in list of all projects
	found := false
	for _, p := range projects {
		if p.ID == project.ID {
			found = true
			break
		}
	}

	if !found {
		t.Fatal("record ID should be found in list of all records")
	}
}
