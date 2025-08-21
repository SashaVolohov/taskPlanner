// A package that defines the application's service.
package service

import (
	"slices"
	"sync"
	"time"

	taskplanner "github.com/SashaVolohov/taskPlanner"
	"github.com/SashaVolohov/taskPlanner/internal/repository"
)

const workersCount = 3

type Service struct {
	repo repository.Task
}

func NewService(repos *repository.Repository) *Service {
	return &Service{repo: repos.Task}
}

func (s *Service) LoadFromFile(path string) error {
	err := s.repo.LoadFromFile(path)
	return err
}

func (s *Service) GetTasksCount() int {
	return s.repo.GetTasksCount()
}

func (s *Service) taskWorker(currentTimeParameters taskplanner.TaskTimeParameters, tasks <-chan taskplanner.TaskInterface, errChannel chan<- error, wg *sync.WaitGroup) {

	defer func() {
		wg.Done()
	}()

Exit:
	for task := range tasks {
		taskTimeParameters := task.GetTaskTimeParameters(currentTimeParameters)
		for i := range taskplanner.CountOfTaskParameters {

			similarFound := slices.Contains(taskTimeParameters[i], currentTimeParameters[i][taskplanner.FirstTimeArgument])
			if !similarFound {
				continue Exit
			}

		}

		task.ExecuteTask(errChannel)
	}
}

func (s *Service) RunTasksByTime(time time.Time, errChannel chan<- error) {

	var currentTimeParameters = taskplanner.TaskTimeParameters{
		[]int{time.Minute()},
		[]int{time.Hour()},
		[]int{time.Day()},
		[]int{int(time.Month())},
		[]int{int(time.Weekday())},
	}

	tasks := make(chan taskplanner.TaskInterface, s.GetTasksCount())

	var wg sync.WaitGroup
	for range workersCount {
		wg.Add(1)
		go s.taskWorker(currentTimeParameters, tasks, errChannel, &wg)
	}

	for _, task := range s.repo.GetTasks() {
		tasks <- task
	}
	close(tasks)

	wg.Wait()

}
