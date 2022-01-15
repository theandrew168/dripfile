package core

type Project struct {
	// readonly (from database, after creation)
	ID int64

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
	Read(id int64) (Project, error)
	Update(project Project) error
	Delete(project Project) error
}
