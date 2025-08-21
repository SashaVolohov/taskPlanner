// The only "repository". Reads information about current tasks from a file
package repository

import (
	"bufio"
	"bytes"
	"os"
	"strconv"
	"strings"

	taskplanner "github.com/SashaVolohov/taskPlanner"
	"github.com/spf13/viper"
)

type TaskRepository struct {
	tasks []taskplanner.TaskInterface
}

func NewTaskRepository() *TaskRepository {
	return &TaskRepository{}
}

func (r *TaskRepository) loadFile(path string) (string, error) {
	f, err := os.Open(path)
	if err != nil {
		return "", err
	}
	defer f.Close()

	buffer := bytes.Buffer{}
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		buffer.WriteString(scanner.Text())
	}

	return buffer.String(), nil
}

func (r *TaskRepository) decodeTaskArgument(i int, argumentParameter string, taskNumericArguments *taskplanner.TaskTimeParameters) error {
	var finalParameter int
	var err error

	if argumentParameter == viper.GetString("anyTimeSymbol") {
		finalParameter = taskplanner.AnyTime
	} else if string(argumentParameter[taskplanner.FirstSymbol]) == viper.GetString("eachSymbol") {
		finalParameter, err = strconv.Atoi(argumentParameter[taskplanner.EachInteger:])
		finalParameter = -finalParameter
	} else {
		finalParameter, err = strconv.Atoi(argumentParameter)
	}

	if err != nil {
		return err
	}

	taskNumericArguments[i] = append(taskNumericArguments[i], finalParameter)
	return nil
}

func (r *TaskRepository) decodeTaskArguments(taskArguments []string, taskNumericArguments *taskplanner.TaskTimeParameters) error {
	for i := range len(taskNumericArguments) {
		argumentParameters := strings.Split(taskArguments[i], viper.GetString("multiTimeSeparationSymbol"))

		for _, argumentParameter := range argumentParameters {
			r.decodeTaskArgument(i, argumentParameter, taskNumericArguments)
		}
	}

	return nil

}

func (r *TaskRepository) decodeTaskString(task string) error {
	taskArguments := strings.Fields(task)
	var taskNumericArguments taskplanner.TaskTimeParameters

	r.decodeTaskArguments(taskArguments, &taskNumericArguments)

	r.tasks = append(r.tasks, taskplanner.NewTask(
		taskNumericArguments[taskplanner.TaskMinute],
		taskNumericArguments[taskplanner.TaskHour],
		taskNumericArguments[taskplanner.TaskDay],
		taskNumericArguments[taskplanner.TaskMonth],
		taskNumericArguments[taskplanner.TaskDayOfWeek],
		strings.Join(taskArguments[taskplanner.TaskCommand:], " ")))

	return nil

}

func (r *TaskRepository) LoadFromFile(path string) (err error) {

	buffer, err := r.loadFile(path)
	if err != nil {
		return err
	}

	tasksList := strings.Split(buffer, viper.GetString("tasksSeparationSymbol"))
	for _, task := range tasksList {
		err = r.decodeTaskString(task)
		if err != nil {
			return err
		}
	}

	return nil

}

func (r *TaskRepository) GetTasksCount() int {
	return len(r.tasks)
}

func (r *TaskRepository) GetTasks() []taskplanner.TaskInterface {
	tasksSliceCopy := make([]taskplanner.TaskInterface, r.GetTasksCount())
	copy(tasksSliceCopy, r.tasks)
	return tasksSliceCopy
}
