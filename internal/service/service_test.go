package service

import (
	"context"
	"os"
	"testing"
	"time"

	taskplanner "github.com/SashaVolohov/taskPlanner"
	"github.com/SashaVolohov/taskPlanner/internal/repository"
)

var testTasksDescription = []struct {
	Minute    []int
	Hour      []int
	Day       []int
	Month     []int
	DayOfWeek []int
	Command   string
}{
	{[]int{0}, []int{0}, []int{10}, []int{10}, []int{taskplanner.AnyTime}, "mkdir test_folder"},
	{[]int{1}, []int{0}, []int{10}, []int{10}, []int{taskplanner.AnyTime}, "mkdir check_engine"},
	{[]int{0}, []int{1}, []int{10}, []int{10}, []int{taskplanner.AnyTime}, "mkdir tu-154"},
}

type FileRepository struct{}

func NewFileRepository() *FileRepository {
	return &FileRepository{}
}

func (s *FileRepository) LoadFromFile(path string) (tasks []taskplanner.Task, err error) {

	for _, task := range testTasksDescription {
		tasks = append(tasks, *taskplanner.NewTask(task.Minute, task.Hour, task.Day, task.Month, task.DayOfWeek, task.Command))
	}

	return tasks, nil
}

func testNewRepository() *repository.Repository {
	return &repository.Repository{
		File: NewFileRepository(),
	}
}

func isDirExists(path string) bool {
	info, err := os.Stat(path)
	if os.IsNotExist(err) {
		return false
	}
	return info.IsDir()
}

func TestMain(t *testing.T) {

	repos := testNewRepository()
	services := NewService(repos)

	ctx, cancel := context.WithCancel(context.Background())

	errChannel := make(chan error)
	go checkErrors(ctx, errChannel, t)

	err := services.LoadFromFile("")
	if err != nil {
		t.Errorf("LoadFromFile failed, test failed.")
	}

	for _, task := range testTasksDescription {
		services.RunTasksByTime(time.Date(time.Now().Year(), time.Month(task.Month[0]), task.Day[0], task.Hour[0], task.Minute[0], 0, 0, time.Now().Location()), errChannel)
	}

	if !isDirExists("test_folder") || !isDirExists("check_engine") || !isDirExists("tu-154") {
		t.Errorf("Folders do not exists, test failed.")
	}

	cancel()

}

func checkErrors(ctx context.Context, out <-chan error, t *testing.T) {
	for {
		select {
		case <-ctx.Done():
			return
		case err := <-out:
			t.Errorf("Error has occurred during task running: %s", err.Error())
		}
	}
}
