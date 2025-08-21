package service

import (
	"testing"

	taskplanner "github.com/SashaVolohov/taskPlanner"
	"github.com/SashaVolohov/taskPlanner/internal/repository"
	"go.uber.org/mock/gomock"
)

func NewRepository(ctrl *gomock.Controller) *repository.Repository {

	task := repository.NewMockTask(ctrl)

	return &repository.Repository{
		Task: task,
	}
}

func TestMain(t *testing.T) {

	ctrl := gomock.NewController(t)
	repos := NewRepository(ctrl)
	services := NewService(repos)

	var testTasksDescription = []struct {
		Minute    []int
		Hour      []int
		Day       []int
		Month     []int
		DayOfWeek []int
		Command   string
		Name      string
		testFunc  func(*Service, *testing.T) bool
	}{
		{[]int{0}, []int{0}, []int{10}, []int{10}, []int{taskplanner.AnyTime}, "1", "First Test", func(services *Service, t *testing.T) bool {
			err := services.LoadFromFile("")
			if err != nil {
				t.Errorf("LoadFromFile failed, test failed.")
			}

			return true

		}},
		{[]int{1}, []int{0}, []int{10}, []int{10}, []int{taskplanner.AnyTime}, "2", "Second Test", func(services *Service, t *testing.T) bool { return true }},
		{[]int{0}, []int{1}, []int{10}, []int{10}, []int{taskplanner.AnyTime}, "3", "Third Test", func(services *Service, t *testing.T) bool { return true }},
	}

	for _, task := range testTasksDescription {

		t.Run(task.Name, func(t *testing.T) {
			t.Parallel()
			if task.testFunc(services, t) == false {
				t.Fatalf(task.Name + " test failed.")
			}
		})
	}

}
