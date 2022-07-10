package storage_test

import (
	"errors"
	"testing"

	"github.com/theandrew168/dripfile/internal/core"
	"github.com/theandrew168/dripfile/internal/database"
	"github.com/theandrew168/dripfile/internal/test"
)

func mockProject() core.Project {
	project := core.NewProject()
	return project
}

func TestProject(t *testing.T) {
	store, closer := test.Storage(t)
	defer closer()

	// create
	project := mockProject()
	err := store.Project.Create(&project)
	if err != nil {
		t.Fatal(err)
	}

	if project.ID == "" {
		t.Fatal("record ID should be non-empty after create")
	}

	// read
	got, err := store.Project.Read(project.ID)
	if err != nil {
		t.Fatal(err)
	}

	if got.ID != project.ID {
		t.Fatal("record ID should be match after read")
	}

	// read all
	projects, err := store.Project.ReadAll()
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

	// delete
	err = store.Project.Delete(project)
	if err != nil {
		t.Fatal(err)
	}

	// verify that ID isn't present anymore
	_, err = store.Project.Read(project.ID)
	if !errors.Is(err, database.ErrNotExist) {
		t.Fatal("record ID should be gone after delete")
	}
}
