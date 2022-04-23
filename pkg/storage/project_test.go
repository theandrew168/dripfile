package storage_test

import (
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
		t.Fatal("project_id should be non-empty after create")
	}
}
