// The package implements a call to task processing logic every second. Very simple...
package handler

import (
	"context"
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

func (h *Handler) ProcessTasks(ctx context.Context) {

	tasksFile := viper.GetString("taskFile")
	logrus.Printf("Loading tasks from file %s...", tasksFile)

	err := h.services.LoadFromFile(tasksFile)
	if err != nil {
		logrus.Fatalf("Failed to read task list from file, check its markup!")
	}

	logrus.Printf("Task scheduler started, %d tasks loaded!", h.services.GetTasksCount())

	errChannel := make(chan error)
	go checkErrors(ctx, errChannel)

	for {
		select {
		case <-ctx.Done():
			return
		case <-time.After(time.Second):
			go h.services.RunTasksByTime(time.Now(), errChannel)
		}
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
