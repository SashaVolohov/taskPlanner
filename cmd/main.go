// The main program package initiates the configuration and application and starts monitoring the task list
package main

import (
	"github.com/SashaVolohov/taskPlanner/internal/handler"
	"github.com/SashaVolohov/taskPlanner/internal/repository"
	"github.com/SashaVolohov/taskPlanner/internal/service"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

func main() {

	logrus.SetFormatter(new(logrus.JSONFormatter))

	if err := initConfig(); err != nil {
		logrus.Fatalf("Error initializing configs: %s", err.Error())
	}

	repos := repository.NewRepository()
	services := service.NewService(repos)
	handlers := handler.NewHandler(services)

	handlers.ProcessTasks()

}

func initConfig() error {
	viper.AddConfigPath("configs")
	viper.SetConfigName("config")
	return viper.ReadInConfig()
}
