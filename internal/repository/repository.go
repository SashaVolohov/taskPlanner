// The package implements access to external resources where information is stored.
// Here is only a file with prescribed tasks
package repository

import taskplanner "github.com/SashaVolohov/taskPlanner"

type Task interface {
	LoadFromFile(path string) (err error)
	GetTasksCount() int
	GetTasks() []taskplanner.TaskInterface
}

type Repository struct {
	Task
}

func NewRepository() *Repository {
	return &Repository{
		Task: NewTaskRepository(),
	}
}
