package services

import (
	"main/utils/logging"

	"go.mongodb.org/mongo-driver/mongo"
)

type Services struct {
	Users *Users
	Tasks *Tasks
}

func NewServices(db *mongo.Database, logger *logging.Logger) *Services {
	usersService := NewUsersService(db, logger)
	tasksService := NewTasksService(db, logger)

	return &Services{
		Users: usersService,
		Tasks: tasksService,
	}
}
