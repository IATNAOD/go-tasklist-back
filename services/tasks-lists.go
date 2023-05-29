package services

import (
	"context"
	"fmt"
	"main/models"
	"main/utils/logging"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type TasksLists struct {
	collection *mongo.Collection
	logger     *logging.Logger
}

func NewTasksListsService(db *mongo.Database, logger *logging.Logger) *TasksLists {
	tasksCollection := db.Collection("tasks-lists")

	return &TasksLists{
		collection: tasksCollection,
		logger:     logger,
	}
}

func (s TasksLists) GetAllUserTasksLists(ctx context.Context, uid string) (tasksLists []models.TasksList, err error) {
	result, err := s.collection.Find(ctx, bson.M{"user_id": uid})
	if result.Err() != nil {
		return tasksLists, err
	}

	err = result.All(ctx, &tasksLists)

	return tasksLists, err
}

func (s TasksLists) GetUserTasksList(ctx context.Context, tlid string, uid string) (tasksList models.TasksList, err error) {
	tloid, err := primitive.ObjectIDFromHex(tlid)
	if err != nil {
		return tasksList, err
	}

	result := s.collection.FindOne(ctx, bson.M{"_id": tloid, "user_id": uid})
	if result.Err() != nil {
		return tasksList, result.Err()
	}

	err = result.Decode(&tasksList)

	return tasksList, err
}

func (s TasksLists) AddTasksList(ctx context.Context, task *models.CreateTasksListDTO) (u models.TasksList, err error) {
	result, err := s.collection.InsertOne(ctx, task)
	if err != nil {
		return u, err
	}

	oid, ok := result.InsertedID.(primitive.ObjectID)
	if !ok {
		return u, fmt.Errorf("failed convert objectid to hex")
	}

	return *task.Build(oid.Hex()), nil
}

func (s TasksLists) UpdateTasksList(ctx context.Context, tlid string, uid string, taskList *models.UpdateTasksListDTO) (t *models.TasksList, err error) {
	tloid, err := primitive.ObjectIDFromHex(tlid)
	if err != nil {
		return t, err
	}

	result := s.collection.FindOneAndUpdate(
		ctx, bson.M{"_id": tloid, "user_id": uid}, bson.M{"$set": taskList},
		options.FindOneAndUpdate().SetReturnDocument(options.After),
	)
	if result.Err() != nil {
		return t, result.Err()
	}

	err = result.Decode(&t)

	return t, err
}

func (s TasksLists) DeleteTasksList(ctx context.Context, tlid string, uid string) (id string, err error) {
	tloid, err := primitive.ObjectIDFromHex(tlid)
	if err != nil {
		return tlid, err
	}

	result, err := s.collection.DeleteOne(ctx, bson.M{"_id": tloid, "user_id": uid})
	if err != nil {
		return tlid, err
	}

	if result.DeletedCount == 0 {
		return "", fmt.Errorf("tasks list not found")
	}

	return tlid, err
}
