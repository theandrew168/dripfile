package storage_test

import (
	"testing"

	"github.com/theandrew168/dripfile/internal/model"
	"github.com/theandrew168/dripfile/internal/postgresql"
	"github.com/theandrew168/dripfile/internal/storage"
	"github.com/theandrew168/dripfile/internal/test"
)

func mockProject() model.Project {
	project := model.NewProject()
	return project
}

func createProject(t *testing.T, store *storage.Storage) (model.Project, DeleterFunc) {
	t.Helper()

	project := mockProject()
	err := store.Project.Create(&project)
	test.AssertNilError(t, err)

	deleter := func(t *testing.T) {
		err := store.Project.Delete(project)
		test.AssertNilError(t, err)
	}

	return project, deleter
}

func TestProjectCreate(t *testing.T) {
	store, closer := test.Storage(t)
	defer closer()

	project, deleter := createProject(t, store)
	defer deleter(t)

	test.AssertNotEqual(t, project.ID, "")
}

func TestProjectDelete(t *testing.T) {
	store, closer := test.Storage(t)
	defer closer()

	project, deleter := createProject(t, store)
	deleter(t)

	_, err := store.Project.Read(project.ID)
	test.AssertErrorIs(t, err, postgresql.ErrNotExist)
}

func TestProjectRead(t *testing.T) {
	store, closer := test.Storage(t)
	defer closer()

	project, deleter := createProject(t, store)
	defer deleter(t)

	got, err := store.Project.Read(project.ID)
	test.AssertNilError(t, err)

	test.AssertEqual(t, got.ID, project.ID)
}

func TestProjectReadAll(t *testing.T) {
	store, closer := test.Storage(t)
	defer closer()

	project1, deleter1 := createProject(t, store)
	defer deleter1(t)

	project2, deleter2 := createProject(t, store)
	defer deleter2(t)

	projects, err := store.Project.ReadAll()
	test.AssertNilError(t, err)

	var ids []string
	for _, r := range projects {
		ids = append(ids, r.ID)
	}

	test.AssertSliceContains(t, ids, project1.ID)
	test.AssertSliceContains(t, ids, project2.ID)
}
