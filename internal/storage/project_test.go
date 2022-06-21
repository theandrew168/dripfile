package storage_test

import (
	"errors"
	"testing"

	"github.com/theandrew168/dripfile/internal/core"
	"github.com/theandrew168/dripfile/internal/database"
	"github.com/theandrew168/dripfile/internal/random"
	"github.com/theandrew168/dripfile/internal/test"
)

func mockProject() core.Project {
	project := core.NewProject(
		random.String(8),
	)
	return project
}

func TestProject(t *testing.T) {
	storage, closer := test.Storage(t)
	defer closer()

	// create
	project := mockProject()
	err := storage.Project.Create(&project)
	if err != nil {
		t.Fatal(err)
	}

	if project.ID == "" {
		t.Fatal("record ID should be non-empty after create")
	}

	// duplicate
	err = storage.Project.Create(&project)
	if !errors.Is(err, database.ErrExist) {
		t.Fatal("duplicate record should return an error")
	}

	// read
	got, err := storage.Project.Read(project.ID)
	if err != nil {
		t.Fatal(err)
	}

	if got.ID != project.ID {
		t.Fatal("record ID should be match after read")
	}

	// read all
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

	// update
	project.CustomerID = random.String(8)
	project.SubscriptionItemID = random.String(8)

	err = storage.Project.Update(project)
	if err != nil {
		t.Fatal(err)
	}

	got, err = storage.Project.Read(project.ID)
	if err != nil {
		t.Fatal(err)
	}

	if got.CustomerID != project.CustomerID {
		t.Fatal("CustomerID should be match after update")
	}
	if got.SubscriptionItemID != project.SubscriptionItemID {
		t.Fatal("SubscriptionItemID should be match after update")
	}

	// delete
	err = storage.Project.Delete(project)
	if err != nil {
		t.Fatal(err)
	}

	// verify that ID isn't present anymore
	_, err = storage.Project.Read(project.ID)
	if !errors.Is(err, database.ErrNotExist) {
		t.Fatal("record ID should be gone after delete")
	}
}
