// The only "repository". Reads information about current tasks from a file
package repository

import (
	"bufio"
	"bytes"
	"os"
)

type TaskListRepository struct{}

func NewTaskListRepository() *TaskListRepository {
	return &TaskListRepository{}
}

func (s *TaskListRepository) LoadFromFile(path string) (string, error) {

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
