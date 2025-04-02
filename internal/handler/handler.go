// The package implements a call to task processing logic every second. Very simple...
package handler

import (
	"time"

	"github.com/SashaVolohov/taskPlanner/internal/service"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

type Handler struct {
	services *service.Service
}

func NewHandler(services *service.Service) *Handler {
	return &Handler{services: services}
}

func (h *Handler) ProcessTasks() {

	tasksFile := viper.GetString("taskFile")
	logrus.Printf("Loading tasks from file %s...", tasksFile)

	err := h.services.TaskList.LoadFromFile(tasksFile)
	if err != nil {
		logrus.Fatalf("Failed to read task list from file, check its markup!")
	}

	logrus.Printf("Task scheduler started, %d tasks loaded!", h.services.TaskList.GetTasksCount())

	errChannel := make(chan error)
	go checkErrors(errChannel)

	for {
		go h.services.TaskList.RunTasksByTime(time.Now(), errChannel)
		time.Sleep(time.Second)
	}
}

func checkErrors(out <-chan error) {
	for {
		err := <-out
		logrus.Printf("Error has occurred during task running: %s", err.Error())
	}
}
