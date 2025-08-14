// A package that defines the application's service.
package service

import (
	"fmt"
	"slices"
	"time"

	taskplanner "github.com/SashaVolohov/taskPlanner"
	"github.com/SashaVolohov/taskPlanner/internal/repository"
)

const workersCount = 3

type Service struct {
	repo  repository.File
	tasks []taskplanner.Task
}

func NewService(repos *repository.Repository) *Service {
	return &Service{repo: repos.File}
}

func (s *Service) LoadFromFile(path string) error {

	var err error
	s.tasks, err = s.repo.LoadFromFile(path)
	return err

}

func (s *Service) GetTasksCount() int {
	return len(s.tasks)
}

func (s *Service) taskWorker(currentTimeParameters taskplanner.TaskTimeParameters, tasks <-chan taskplanner.Task, errChannel chan<- error) {
Exit:
	for task := range tasks {
		taskTimeParameters := task.GetTaskTimeParameters(currentTimeParameters)
		for i := range taskplanner.CountOfTaskParameters {

			similarFound := slices.Contains(taskTimeParameters[i], currentTimeParameters[i][taskplanner.FirstTimeArgument])
			if !similarFound {
				fmt.Print(taskTimeParameters, " - ", currentTimeParameters)
				continue Exit
			}

		}

		task.ExecuteTask(errChannel)
	}
}

func (s *Service) RunTasksByTime(time time.Time, errChannel chan<- error) {

	fmt.Print(time.GoString())

	var currentTimeParameters taskplanner.TaskTimeParameters

	currentTimeParameters[taskplanner.TaskMinute] = append(currentTimeParameters[taskplanner.TaskMinute], time.Minute())
	currentTimeParameters[taskplanner.TaskHour] = append(currentTimeParameters[taskplanner.TaskHour], time.Hour())
	currentTimeParameters[taskplanner.TaskDay] = append(currentTimeParameters[taskplanner.TaskDay], time.Day())
	currentTimeParameters[taskplanner.TaskMonth] = append(currentTimeParameters[taskplanner.TaskMonth], int(time.Month()))
	currentTimeParameters[taskplanner.TaskDayOfWeek] = append(currentTimeParameters[taskplanner.TaskDayOfWeek], int(time.Weekday()))

	tasks := make(chan taskplanner.Task, len(s.tasks))

	for range workersCount {
		go s.taskWorker(currentTimeParameters, tasks, errChannel)
	}

	for _, task := range s.tasks {
		tasks <- task
	}
	close(tasks)

}
