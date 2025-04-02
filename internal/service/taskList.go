// The main logic of the program is here.
// Creates task entities when the program starts, and then processes them on command of the handler
package service

import (
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"

	"slices"

	taskplanner "github.com/SashaVolohov/taskPlanner"
	"github.com/SashaVolohov/taskPlanner/internal/repository"
	"github.com/spf13/viper"
)

const (
	programName = iota
	firstArgument
)

type TaskListService struct {
	repo  repository.TaskList
	tasks []taskplanner.Task
}

func NewTaskListService(repo repository.TaskList) *TaskListService {
	return &TaskListService{repo: repo}
}

func (s *TaskListService) LoadFromFile(path string) error {

	buffer, err := s.repo.LoadFromFile(path)
	if err != nil {
		return err
	}

	tasksStrings := strings.Split(buffer, viper.GetString("tasksSeparationSymbol"))
	for _, taskString := range tasksStrings {
		stringArguments := strings.Fields(taskString)
		var numberArguments taskplanner.TaskTimeParameters

		for i := taskplanner.TaskMinute; i < len(numberArguments); i++ {

			integers := strings.Split(stringArguments[i], viper.GetString("multiTimeSeparationSymbol"))

			for _, v := range integers {
				integer, err := strconv.Atoi(v)

				if err != nil {
					if v == viper.GetString("anyTimeSymbol") {
						numberArguments[i] = append(numberArguments[i], taskplanner.AnyTime)
						break
					} else if string(v[taskplanner.FirstSymbol]) == viper.GetString("eachSymbol") {

						integer, err = strconv.Atoi(v[taskplanner.EachInteger:])
						integer = -integer
						if err != nil {
							return err
						}

					} else {
						return err
					}
				}

				numberArguments[i] = append(numberArguments[i], integer)
			}

		}

		s.tasks = append(s.tasks, *taskplanner.NewTask(
			numberArguments[taskplanner.TaskMinute],
			numberArguments[taskplanner.TaskHour],
			numberArguments[taskplanner.TaskDay],
			numberArguments[taskplanner.TaskMonth],
			numberArguments[taskplanner.TaskDayOfWeek],
			strings.Join(stringArguments[taskplanner.TaskCommand:], " ")))
	}

	return nil

}

func (s *TaskListService) GetTasksCount() int {
	return len(s.tasks)
}

func (s *TaskListService) RunTasksByTime(time time.Time, out chan<- error) {

	var currentTimeParameters taskplanner.TaskTimeParameters

	currentTimeParameters[taskplanner.TaskMinute] = append(currentTimeParameters[taskplanner.TaskMinute], time.Minute())
	currentTimeParameters[taskplanner.TaskHour] = append(currentTimeParameters[taskplanner.TaskHour], time.Hour())
	currentTimeParameters[taskplanner.TaskDay] = append(currentTimeParameters[taskplanner.TaskDay], time.Day())
	currentTimeParameters[taskplanner.TaskMonth] = append(currentTimeParameters[taskplanner.TaskMonth], int(time.Month()))
	currentTimeParameters[taskplanner.TaskDayOfWeek] = append(currentTimeParameters[taskplanner.TaskDayOfWeek], int(time.Weekday()))

	for _, task := range s.tasks {

		taskTimeParameters := task.GetTaskTimeParameters(currentTimeParameters)
		for i := range taskplanner.CountOfNumberArguments {

			similarFound := slices.Contains(taskTimeParameters[i], currentTimeParameters[i][taskplanner.FirstTimeArgument])
			if !similarFound {
				return
			}

		}

		command := strings.Fields(task.GetCommandString())
		cmd := exec.Command(command[programName], command[firstArgument:]...)

		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		cmd.Stdin = os.Stdin

		err := cmd.Run()
		if err != nil {
			out <- err
			return
		}

	}

}
