package storage_test

import (
	"testing"

	"github.com/theandrew168/dripfile/internal/model"
	"github.com/theandrew168/dripfile/internal/postgresql"
	"github.com/theandrew168/dripfile/internal/test"
)

func mockProject() model.Project {
	project := model.NewProject()
	return project
}

func TestProject(t *testing.T) {
	store, closer := test.Storage(t)
	defer closer()

	// create
	project := mockProject()
	err := store.Project.Create(&project)
	test.AssertNilError(t, err)

	test.AssertNotEqual(t, project.ID, "")

	// read
	got, err := store.Project.Read(project.ID)
	test.AssertNilError(t, err)

	test.AssertEqual(t, got.ID, project.ID)

	// read all
	projects, err := store.Project.ReadAll()
	test.AssertNilError(t, err)

	// verify that ID is present in list of all projects
	var ids []string
	for _, r := range projects {
		ids = append(ids, r.ID)
	}

	test.AssertSliceContains(t, ids, project.ID)

	// delete
	err = store.Project.Delete(project)
	test.AssertNilError(t, err)

	// verify that ID isn't present anymore
	_, err = store.Project.Read(project.ID)
	test.AssertErrorIs(t, err, postgresql.ErrNotExist)
}
