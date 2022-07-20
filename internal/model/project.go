package model

type Project struct {
	// readonly (from database, after creation)
	ID string
}

func NewProject() Project {
	project := Project{}
	return project
}
