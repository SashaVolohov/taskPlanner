// A package that defines the application's services. There is only one service - Tasklist
package service

import (
	"time"

	"github.com/SashaVolohov/taskPlanner/internal/repository"
)

type TaskList interface {
	LoadFromFile(path string) error
	GetTasksCount() int
	RunTasksByTime(time time.Time, out chan<- error)
}

type Service struct {
	TaskList
}

func NewService(repos *repository.Repository) *Service {
	return &Service{
		TaskList: NewTaskListService(repos.TaskList),
	}
}
