package models

import "time"

type TasksList struct {
	ID        string    `json:"id" bson:"_id,omitempty"`
	UserID    string    `json:"user_id" bson:"user_id" validate:"nonzero,len=24"`
	Name      string    `json:"name" bson:"name" validate:"nonzero"`
	Color     string    `json:"color" bson:"color" validate:"nonzero"`
	Hidden    bool      `json:"hidden" bson:"hidden" validate:"nonnil"`
	UpdatedAt time.Time `json:"UpdatedAt" bson:"UpdatedAt" validate:"nonzero"`
	CreatedAt time.Time `json:"CreatedAt" bson:"CreatedAt" validate:"nonzero"`
}

type CreateTasksListRB struct {
	Name   string `json:"name" bson:"name" validate:"nonzero"`
	Color  string `json:"color" bson:"color" validate:"nonzero"`
	Hidden bool   `json:"hidden" bson:"hidden" validate:"nonnil"`
}

type CreateTasksListDTO struct {
	UserID    string    `json:"user_id" bson:"user_id" validate:"nonzero,len=24"`
	Name      string    `json:"name" bson:"name" validate:"nonzero"`
	Color     string    `json:"color" bson:"color" validate:"nonzero"`
	Hidden    bool      `json:"hidden" bson:"hidden" validate:"nonnil"`
	UpdatedAt time.Time `json:"UpdatedAt" bson:"UpdatedAt" validate:"nonzero"`
	CreatedAt time.Time `json:"CreatedAt" bson:"CreatedAt" validate:"nonzero"`
}

type UpdateTasksListRB struct {
	ID     string `json:"id" bson:"_id,omitempty" validate:"nonzero"`
	Name   string `json:"name" bson:"name"`
	Color  string `json:"color" bson:"color"`
	Hidden bool   `json:"hidden" bson:"hidden"`
}

type UpdateTasksListDTO struct {
	Name      string    `json:"name" bson:"name"`
	Color     string    `json:"color" bson:"color"`
	Hidden    bool      `json:"hidden" bson:"hidden"`
	UpdatedAt time.Time `json:"UpdatedAt" bson:"UpdatedAt"`
}

type DeleteTasksListDTO struct {
	ID string `json:"id" bson:"_id,omitempty" validate:"nonzero,len=24"`
}

func (task CreateTasksListRB) Build(uid string) *CreateTasksListDTO {
	return &CreateTasksListDTO{
		UserID:    uid,
		Name:      task.Name,
		Color:     task.Color,
		Hidden:    task.Hidden,
		UpdatedAt: time.Now(),
		CreatedAt: time.Now(),
	}
}

func (task CreateTasksListDTO) Build(id string) *TasksList {
	return &TasksList{
		ID:        id,
		UserID:    task.UserID,
		Name:      task.Name,
		Color:     task.Color,
		Hidden:    task.Hidden,
		UpdatedAt: task.UpdatedAt,
		CreatedAt: task.CreatedAt,
	}
}

func (task UpdateTasksListRB) Build() *UpdateTasksListDTO {
	return &UpdateTasksListDTO{
		Name:      task.Name,
		Color:     task.Color,
		Hidden:    task.Hidden,
		UpdatedAt: time.Now(),
	}
}
