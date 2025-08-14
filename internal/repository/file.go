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

type FileRepository struct{}

func NewFileRepository() *FileRepository {
	return &FileRepository{}
}

func (s *FileRepository) loadFile(path string) (string, error) {
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

func (s *FileRepository) LoadFromFile(path string) (tasks []taskplanner.Task, err error) {

	buffer, err := s.loadFile(path)
	if err != nil {
		return nil, err
	}

	tasksList := strings.Split(buffer, viper.GetString("tasksSeparationSymbol"))
	for _, task := range tasksList {
		taskArguments := strings.Fields(task)
		var taskNumericArguments taskplanner.TaskTimeParameters

		for i := range len(taskNumericArguments) {

			argumentParameters := strings.Split(taskArguments[i], viper.GetString("multiTimeSeparationSymbol"))

			for _, argumentParameter := range argumentParameters {

				var numericParameter int

				if argumentParameter == viper.GetString("anyTimeSymbol") {
					taskNumericArguments[i] = append(taskNumericArguments[i], taskplanner.AnyTime)
					break
				}

				if string(argumentParameter[taskplanner.FirstSymbol]) == viper.GetString("eachSymbol") {

					numericParameter, err := strconv.Atoi(argumentParameter[taskplanner.EachInteger:])
					if err != nil {
						return nil, err
					}

					taskNumericArguments[i] = append(taskNumericArguments[i], -numericParameter)
					break

				}

				numericParameter, err := strconv.Atoi(argumentParameter)
				if err != nil {
					return nil, err
				}

				taskNumericArguments[i] = append(taskNumericArguments[i], numericParameter)
			}

		}

		tasks = append(tasks, *taskplanner.NewTask(
			taskNumericArguments[taskplanner.TaskMinute],
			taskNumericArguments[taskplanner.TaskHour],
			taskNumericArguments[taskplanner.TaskDay],
			taskNumericArguments[taskplanner.TaskMonth],
			taskNumericArguments[taskplanner.TaskDayOfWeek],
			strings.Join(taskArguments[taskplanner.TaskCommand:], " ")))
	}

	return tasks, nil

}
