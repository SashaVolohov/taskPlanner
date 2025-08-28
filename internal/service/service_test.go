package service

import (
	"context"
	"testing"
	"time"

	taskplanner "github.com/SashaVolohov/taskPlanner"
	"github.com/SashaVolohov/taskPlanner/internal/repository"
	"github.com/sirupsen/logrus"
	"go.uber.org/mock/gomock"
)

func TestMain(t *testing.T) {

	ctrl := gomock.NewController(t)
	repo := repository.NewMockTask(ctrl)
	services := NewService(repo)

	ctx, cancel := context.WithCancel(context.Background())

	defer cancel()

	errChannel := make(chan error)
	go checkErrors(ctx, errChannel)

	type args struct {
		Minute    []int
		Hour      []int
		Day       []int
		Month     []int
		DayOfWeek []int
		Command   string
		Name      string
	}

	var tests = []struct {
		args     args
		name     string
		testFunc func(args, *Service, *testing.T) bool
	}{
		{
			args: args{
				Minute:    []int{0},
				Hour:      []int{0},
				Day:       []int{10},
				Month:     []int{10},
				DayOfWeek: []int{taskplanner.AnyTime},
				Command:   "1",
			},
			name: "Execute Task",
			testFunc: func(args args, services *Service, t *testing.T) bool {

				taskTime := time.Date(time.Now().Year(), time.Month(args.Month[0]), args.Day[0], args.Hour[0], args.Minute[0], 0, 0, time.Now().Location())
				var wantedTaskParameters taskplanner.TaskTimeParameters = [taskplanner.CountOfTaskParameters][]int{args.Minute, args.Hour, args.Day, args.Month, {int(taskTime.Weekday())}}

				testTask := taskplanner.NewMockTaskInterface(ctrl)
				testTask.EXPECT().GetTaskTimeParameters(wantedTaskParameters).Return(wantedTaskParameters)

				testTask.EXPECT().ExecuteTask(errChannel).Return()

				repo.EXPECT().LoadFromFile("").Return(nil)
				repo.EXPECT().GetTasksCount().Return(1)
				repo.EXPECT().GetTasks().Return([]taskplanner.TaskInterface{testTask})

				err := services.LoadFromFile("")
				if err != nil {
					t.Errorf("LoadFromFile failed, test failed.")
				}

				services.RunTasksByTime(taskTime, errChannel)

				return true

			}},
		{
			args: args{
				Minute:    []int{1},
				Hour:      []int{0},
				Day:       []int{10},
				Month:     []int{10},
				DayOfWeek: []int{taskplanner.AnyTime},
				Command:   "1",
			},
			name: "Bad Time For Execute",
			testFunc: func(args args, services *Service, t *testing.T) bool {

				taskTime := time.Date(time.Now().Year(), time.Month(args.Month[0]), args.Day[0], args.Hour[0], 0, 0, 0, time.Now().Location())
				var wantedTaskParameters taskplanner.TaskTimeParameters = [taskplanner.CountOfTaskParameters][]int{args.Minute, args.Hour, args.Day, args.Month, {int(taskTime.Weekday())}}

				testTask := taskplanner.NewMockTaskInterface(ctrl)
				testTask.EXPECT().GetTaskTimeParameters([taskplanner.CountOfTaskParameters][]int{{0}, args.Hour, args.Day, args.Month, {int(taskTime.Weekday())}}).
					Return(wantedTaskParameters)

				repo.EXPECT().LoadFromFile("").Return(nil)
				repo.EXPECT().GetTasksCount().Return(1)
				repo.EXPECT().GetTasks().Return([]taskplanner.TaskInterface{testTask})

				err := services.LoadFromFile("")
				if err != nil {
					t.Errorf("LoadFromFile failed, test failed.")
				}

				services.RunTasksByTime(taskTime, errChannel)

				return true

			}},
		{
			args: args{
				Minute:    []int{0, 1},
				Hour:      []int{0, 0},
				Day:       []int{10, 10},
				Month:     []int{10, 10},
				DayOfWeek: []int{taskplanner.AnyTime, taskplanner.AnyTime},
				Command:   "1",
			},
			name: "Two tasks(one should be runned, other - no)",
			testFunc: func(args args, services *Service, t *testing.T) bool {

				taskTime := time.Date(time.Now().Year(), time.Month(args.Month[0]), args.Day[0], args.Hour[0], args.Minute[0], 0, 0, time.Now().Location())
				var wantedTaskParameters taskplanner.TaskTimeParameters = [taskplanner.CountOfTaskParameters][]int{{args.Minute[0]}, {args.Hour[0]}, {args.Day[0]}, {args.Month[0]}, {int(taskTime.Weekday())}}

				tasks := []*taskplanner.MockTaskInterface{
					taskplanner.NewMockTaskInterface(ctrl),
					taskplanner.NewMockTaskInterface(ctrl),
				}

				tasks[0].EXPECT().GetTaskTimeParameters(wantedTaskParameters).Return(wantedTaskParameters)
				tasks[0].EXPECT().ExecuteTask(errChannel).Return()

				tasks[1].EXPECT().GetTaskTimeParameters(wantedTaskParameters).
					Return([taskplanner.CountOfTaskParameters][]int{{args.Minute[1]}, {args.Hour[1]}, {args.Day[1]}, {args.Month[1]}, {int(taskTime.Weekday())}})

				repo.EXPECT().LoadFromFile("").Return(nil)
				repo.EXPECT().GetTasksCount().Return(len(tasks))
				repo.EXPECT().GetTasks().Return([]taskplanner.TaskInterface{tasks[0], tasks[1]})

				err := services.LoadFromFile("")
				if err != nil {
					t.Errorf("LoadFromFile failed, test failed.")
				}

				services.RunTasksByTime(taskTime, errChannel)

				return true

			}},
	}

	for _, task := range tests {

		t.Run(task.name, func(t *testing.T) {
			t.Parallel()
			if task.testFunc(task.args, services, t) == false {
				t.Fatalf(task.name + " test failed.")
			}
		})
	}

}

func checkErrors(ctx context.Context, out <-chan error) {
	for {
		select {
		case <-ctx.Done():
			return
		case err := <-out:
			logrus.Printf("Error has occurred during task running: %s", err.Error())
		}
	}
}
