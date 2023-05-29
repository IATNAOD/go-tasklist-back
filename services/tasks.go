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

type Tasks struct {
	collection *mongo.Collection
	logger     *logging.Logger
}

func NewTasksService(db *mongo.Database, logger *logging.Logger) *Tasks {
	tasksCollection := db.Collection("tasks")

	return &Tasks{
		collection: tasksCollection,
		logger:     logger,
	}
}

func (s Tasks) GetAllUserTasks(ctx context.Context, uid string) (tasks []models.Task, err error) {
	result, err := s.collection.Find(ctx, bson.M{"user_id": uid})
	if result.Err() != nil {
		return tasks, err
	}

	err = result.All(ctx, &tasks)

	return tasks, err
}

func (s Tasks) AddTask(ctx context.Context, task *models.CreateTaskDTO) (u models.Task, err error) {
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

func (s Tasks) UpdateTask(ctx context.Context, tid string, uid string, task *models.UpdateTaskDTO) (t *models.Task, err error) {
	toid, err := primitive.ObjectIDFromHex(tid)
	if err != nil {
		return t, err
	}

	result := s.collection.FindOneAndUpdate(
		ctx, bson.M{"_id": toid, "user_id": uid}, bson.M{"$set": task},
		options.FindOneAndUpdate().SetReturnDocument(options.After),
	)
	if result.Err() != nil {
		return t, result.Err()
	}

	err = result.Decode(&t)

	return t, err
}

func (s Tasks) DeleteTask(ctx context.Context, tid string, uid string) (id string, err error) {
	toid, err := primitive.ObjectIDFromHex(tid)
	if err != nil {
		return tid, err
	}

	result, err := s.collection.DeleteOne(ctx, bson.M{"_id": toid, "user_id": uid})
	if err != nil {
		return tid, err
	}

	if result.DeletedCount == 0 {
		return "", fmt.Errorf("task not found")
	}

	return tid, err
}

func (s Tasks) DeleteAllTask(ctx context.Context, uid string) (err error) {
	result, err := s.collection.DeleteMany(ctx, bson.M{"user_id": uid})
	if err != nil {
		return err
	}

	if result.DeletedCount == 0 {
		return fmt.Errorf("tasks not found")
	}

	return err
}
