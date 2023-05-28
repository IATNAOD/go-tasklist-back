package models

type Task struct {
	ID     string `json:"id" bson:"_id,omitempty"`
	UserID string `json:"user_id" bson:"user_id" validate:"nonzero,len=24"`
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
	Title  string `json:"title" bson:"title" validate:"nonzero"`
	Note   string `json:"note" bson:"note"`
	Subs   []struct {
		Title    string `json:"title" bson:"title"`
		Complete bool   `json:"complete" bson:"complete"`
	} `json:"subs" bson:"subs"`
	Complete bool `json:"complete" bson:"complete" validate:"nonnil"`
}

type UpdateTaskDTO struct {
	Title string `json:"title" bson:"title" validate:"nonzero"`
	Note  string `json:"note" bson:"note"`
	Subs  []struct {
		Title    string `json:"title" bson:"title"`
		Complete bool   `json:"complete" bson:"complete"`
	} `json:"subs" bson:"subs"`
	Complete bool `json:"complete" bson:"complete" validate:"nonnil"`
}

type DeleteTaskDTO struct {
	ID string `json:"id" bson:"_id,omitempty"`
}

func (t Task) BuildForUpdate() *UpdateTaskDTO {
	return &UpdateTaskDTO{
		Title:    t.Title,
		Note:     t.Note,
		Subs:     t.Subs,
		Complete: t.Complete,
	}
}

func (t CreateTaskDTO) Build(id string) *Task {
	return &Task{
		ID:       id,
		UserID:   t.UserID,
		Title:    t.Title,
		Note:     t.Note,
		Subs:     t.Subs,
		Complete: t.Complete,
	}
}
