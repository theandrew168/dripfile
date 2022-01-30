package core

type Project struct {
	// readonly (from database, after creation)
	ID string
}

func NewProject() Project {
	project := Project{}
	return project
}

type ProjectStorage interface {
	Create(project *Project) error
	Read(id string) (Project, error)
	Delete(project Project) error
}
