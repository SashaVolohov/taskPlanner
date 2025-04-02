// The package implements access to external resources where information is stored.
// Here is only a file with prescribed tasks
package repository

type TaskList interface {
	LoadFromFile(path string) (string, error)
}

type Repository struct {
	TaskList
}

func NewRepository() *Repository {
	return &Repository{
		TaskList: NewTaskListRepository(),
	}
}
