package models

import "time"

type Task struct {
	ID     string `json:"id" bson:"_id,omitempty"`
	UserID string `json:"user_id" bson:"user_id" validate:"nonzero,len=24"`
	ListID string `json:"list_id" bson:"list_id" validate:"nonzero,len=24"`
	Title  string `json:"title" bson:"title" validate:"nonzero"`
	Note   string `json:"note" bson:"note"`
	Subs   []struct {
		Title    string `json:"title" bson:"title"`
		Complete bool   `json:"complete" bson:"complete"`
	} `json:"subs" bson:"subs"`
	Complete  bool      `json:"complete" bson:"complete" validate:"nonnil"`
	UpdatedAt time.Time `json:"UpdatedAt" bson:"UpdatedAt" validate:"nonzero"`
	CreatedAt time.Time `json:"CreatedAt" bson:"CreatedAt" validate:"nonzero"`
}

type CreateTaskRB struct {
	ListID string `json:"list_id" bson:"list_id" validate:"nonzero,len=24"`
	Title  string `json:"title" bson:"title" validate:"nonzero"`
	Note   string `json:"note" bson:"note"`
	Subs   []struct {
		Title    string `json:"title" bson:"title"`
		Complete bool   `json:"complete" bson:"complete"`
	} `json:"subs" bson:"subs"`
	Complete bool `json:"complete" bson:"complete" validate:"nonnil"`
}

type CreateTaskDTO struct {
	UserID string `json:"user_id" bson:"user_id" validate:"nonzero,len=24"`
	ListID string `json:"list_id" bson:"list_id" validate:"nonzero,len=24"`
	Title  string `json:"title" bson:"title" validate:"nonzero"`
	Note   string `json:"note" bson:"note"`
	Subs   []struct {
		Title    string `json:"title" bson:"title"`
		Complete bool   `json:"complete" bson:"complete"`
	} `json:"subs" bson:"subs"`
	Complete  bool      `json:"complete" bson:"complete" validate:"nonnil"`
	UpdatedAt time.Time `json:"UpdatedAt" bson:"UpdatedAt" validate:"nonzero"`
	CreatedAt time.Time `json:"CreatedAt" bson:"CreatedAt" validate:"nonzero"`
}

type UpdateTaskRB struct {
	ID     string `json:"id" bson:"_id,omitempty"`
	ListID string `json:"list_id" bson:"list_id" validate:"nonzero,len=24"`
	Title  string `json:"title" bson:"title" validate:"nonzero"`
	Note   string `json:"note" bson:"note"`
	Subs   []struct {
		Title    string `json:"title" bson:"title"`
		Complete bool   `json:"complete" bson:"complete"`
	} `json:"subs" bson:"subs"`
	Complete bool `json:"complete" bson:"complete" validate:"nonnil"`
}

type UpdateTaskDTO struct {
	ListID string `json:"list_id" bson:"list_id" validate:"nonzero,len=24"`
	Title  string `json:"title" bson:"title" validate:"nonzero"`
	Note   string `json:"note" bson:"note"`
	Subs   []struct {
		Title    string `json:"title" bson:"title"`
		Complete bool   `json:"complete" bson:"complete"`
	} `json:"subs" bson:"subs"`
	Complete  bool      `json:"complete" bson:"complete" validate:"nonnil"`
	UpdatedAt time.Time `json:"UpdatedAt" bson:"UpdatedAt" validate:"nonzero"`
}

type DeleteTaskDTO struct {
	ID string `json:"id" bson:"_id,omitempty"`
}

func (t CreateTaskRB) Build(uid string) *CreateTaskDTO {
	return &CreateTaskDTO{
		UserID:    uid,
		ListID:    t.ListID,
		Title:     t.Title,
		Note:      t.Note,
		Subs:      t.Subs,
		Complete:  t.Complete,
		UpdatedAt: time.Now(),
		CreatedAt: time.Now(),
	}
}

func (t UpdateTaskRB) Build() *UpdateTaskDTO {
	return &UpdateTaskDTO{
		ListID:    t.ListID,
		Title:     t.Title,
		Note:      t.Note,
		Subs:      t.Subs,
		Complete:  t.Complete,
		UpdatedAt: time.Now(),
	}
}

func (t CreateTaskDTO) Build(id string) *Task {
	return &Task{
		ID:        id,
		UserID:    t.UserID,
		ListID:    t.ListID,
		Title:     t.Title,
		Note:      t.Note,
		Subs:      t.Subs,
		Complete:  t.Complete,
		UpdatedAt: t.UpdatedAt,
		CreatedAt: t.CreatedAt,
	}
}
