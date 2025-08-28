package repository

import (
	"slices"
	"testing"

	taskplanner "github.com/SashaVolohov/taskPlanner"
	"github.com/spf13/viper"
)

func isArrayEquals(a taskplanner.TaskTimeParameters, b taskplanner.TaskTimeParameters) bool {
	if len(a) != len(b) {
		return false
	}
	for i, v := range a {
		if !slices.Equal(v, b[i]) {
			return false
		}
	}
	return true
}

func initConfig() error {
	viper.AddConfigPath("../../configs")
	viper.SetConfigName("config")
	return viper.ReadInConfig()
}

func TestMain(t *testing.T) {

	err := initConfig()
	if err != nil {
		t.Errorf("Failed to load config, test failed - %s", err.Error())
	}

	type args struct {
		Minute    []int
		Hour      []int
		Day       []int
		Month     []int
		DayOfWeek []int
	}

	var tests = []struct {
		args     args
		name     string
		testFunc func(args, *testing.T) bool
	}{
		{
			args: args{
				Minute:    []int{0, 0, 0},
				Hour:      []int{0, 10, 18},
				Day:       []int{10, 10, 1},
				Month:     []int{10, 10, 10},
				DayOfWeek: []int{0, 0, 0},
			},
			name: "Load 3 tasks",
			testFunc: func(args args, t *testing.T) bool {

				const tasksCount = 3

				repository := NewRepository()
				err = repository.LoadFromFile("../../test/repositoryTestTasks.txt")
				if err != nil {
					t.Fatalf("Failed to read file - %s", err.Error())
					return false
				}

				if repository.GetTasksCount() != tasksCount {
					t.Errorf("Incorrect tasks count - want 3, but got %d", repository.GetTasksCount())
					return false
				}

				tasks := repository.GetTasks()

				for i := range tasksCount {
					wantTaskParameters := taskplanner.TaskTimeParameters{{args.Minute[i]}, {args.Hour[i]}, {args.Day[i]}, {args.Month[i]}, {args.DayOfWeek[i]}}
					timeParameters := tasks[i].GetTaskTimeParameters(wantTaskParameters)
					if !isArrayEquals(timeParameters, wantTaskParameters) {
						t.Errorf("Incorrect tasks time")
						return false
					}
				}

				return true

			}},
	}

	for _, task := range tests {

		t.Run(task.name, func(t *testing.T) {
			t.Parallel()
			if task.testFunc(task.args, t) == false {
				t.Fatalf(task.name + " test failed.")
			}
		})
	}

}
