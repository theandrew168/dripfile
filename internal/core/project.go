package core

type Project struct {
	// readonly (from database, after creation)
	ID string

	Name string
}

func NewProject(name string) Project {
	project := Project{
		Name: name,
	}
	return project
}

type ProjectStorage interface {
	Create(project *Project) error
	Read(id string) (Project, error)
	Update(project Project) error
	Delete(project Project) error
}
