package core

import (
	"context"
)

type Project struct {
	Name string

	// readonly (from database, after creation)
	ID int
}

func NewProject(name string) Project {
	project := Project{
		Name: name,
	}
	return project
}

type ProjectStorage interface {
	CreateProject(ctx context.Context, project *Project) error
	ReadProject(ctx context.Context, id int) (Project, error)
	UpdateProject(ctx context.Context, project Project) error
	DeleteProject(ctx context.Context, project Project) error
}
