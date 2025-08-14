// The package implements access to external resources where information is stored.
// Here is only a file with prescribed tasks
package repository

import taskplanner "github.com/SashaVolohov/taskPlanner"

type File interface {
	LoadFromFile(path string) (tasks []taskplanner.Task, err error)
}

type Repository struct {
	File
}

func NewRepository() *Repository {
	return &Repository{
		File: NewFileRepository(),
	}
}
