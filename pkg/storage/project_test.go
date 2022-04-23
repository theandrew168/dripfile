package storage_test

import (
	"errors"
	"testing"

	"github.com/theandrew168/dripfile/pkg/core"
	"github.com/theandrew168/dripfile/pkg/database"
	"github.com/theandrew168/dripfile/pkg/storage"
	"github.com/theandrew168/dripfile/pkg/test"
)

func TestCreateProject(t *testing.T) {
	cfg := test.Config(t)

	pool, err := database.ConnectPool(cfg.DatabaseURI)
	if err != nil {
		t.Fatal(err)
	}
	defer pool.Close()

	storage := storage.New(pool)
	project := core.NewMockProject()
	err = storage.Project.Create(&project)
	if err != nil {
		t.Fatal(err)
	}

	if project.ID == "" {
		t.Fatal("record ID should be non-empty after create")
	}
}

func TestCreateProjectDuplicate(t *testing.T) {
	cfg := test.Config(t)

	pool, err := database.ConnectPool(cfg.DatabaseURI)
	if err != nil {
		t.Fatal(err)
	}
	defer pool.Close()

	storage := storage.New(pool)
	project := core.NewMockProject()
	err = storage.Project.Create(&project)
	if err != nil {
		t.Fatal(err)
	}

	// attempt to create the same project again
	err = storage.Project.Create(&project)
	if !errors.Is(err, core.ErrExist) {
		t.Fatal("duplicate record should return an error")
	}
}
