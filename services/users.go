package services

import (
	"context"
	"fmt"
	"main/models"
	"main/utils/logging"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type Users struct {
	collection *mongo.Collection
	logger     *logging.Logger
}

func NewUsersService(db *mongo.Database, logger *logging.Logger) *Users {
	usersCollection := db.Collection("users")

	return &Users{
		collection: usersCollection,
		logger:     logger,
	}
}

func (s Users) CreateUser(ctx context.Context, user *models.User) (string, error) {
	result, err := s.collection.InsertOne(ctx, user)
	if err != nil {
		return "", err
	}

	oid, ok := result.InsertedID.(primitive.ObjectID)
	if !ok {
		return "", fmt.Errorf("failed convert objectid to hex")
	}

	return oid.Hex(), nil
}

func (s Users) FindUserByEmail(ctx context.Context, email string) (u models.User, err error) {
	result := s.collection.FindOne(ctx, bson.M{"email": email})
	if result.Err() != nil {
		return u, result.Err()
	}

	err = result.Decode(&u)

	return u, err
}

func (s Users) GetAllUsers(ctx context.Context) (users []models.User, err error) {
	result, err := s.collection.Find(ctx, bson.M{})
	if result.Err() != nil {
		return users, err
	}

	err = result.All(ctx, &users)

	return users, err
}
