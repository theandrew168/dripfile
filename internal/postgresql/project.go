package postgresql

import (
	"github.com/jackc/pgx/v4/pgxpool"

	"github.com/theandrew168/dripfile/internal/core"
)

type projectStorage struct {
	conn *pgxpool.Pool
}

func NewProjectStorage(conn *pgxpool.Pool) core.ProjectStorage {
	s := projectStorage{
		conn: conn,
	}
	return &s
}

func (s *projectStorage) Create(project *core.Project) error {
	return nil
}

func (s *projectStorage) Read(id int64) (core.Project, error) {
	return core.Project{}, nil
}

func (s *projectStorage) Update(project core.Project) error {
	return nil
}

func (s *projectStorage) Delete(project core.Project) error {
	return nil
}
