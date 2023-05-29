package routes

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"main/middlewares"
	"main/models"
	"main/services"
	"main/utils/logging"
	"net/http"

	"github.com/julienschmidt/httprouter"
	"github.com/redis/go-redis/v9"
	"gopkg.in/validator.v2"
)

type TasksListsHandler struct {
	Parent      *Router
	Router      *httprouter.Router
	Services    *services.Services
	middlewares *middlewares.Middlewares
	logger      *logging.Logger
	redis       *redis.Client
}

func NewTasksListsHandler(router *Router) *TasksListsHandler {
	return &TasksListsHandler{
		Parent:      router,
		Router:      router.Router,
		Services:    router.Services,
		middlewares: router.middlewares,
		logger:      router.logger,
		redis:       router.redis,
	}
}

func (h TasksListsHandler) RegisterTasksListsRoutes() {
	h.Router.HandlerFunc(http.MethodGet, "/tasks-lists/", h.middlewares.ApplyMiddlewares(
		h.GetAllTasksLists,
		h.middlewares.ForAuth,
	))
	h.Router.HandlerFunc(http.MethodPost, "/tasks-lists/", h.middlewares.ApplyMiddlewares(
		h.AddNewTasksList,
		h.middlewares.ForAuth,
	))
	h.Router.HandlerFunc(http.MethodPatch, "/tasks-lists/", h.middlewares.ApplyMiddlewares(
		h.UpdateTasksList,
		h.middlewares.ForAuth,
	))
	h.Router.HandlerFunc(http.MethodDelete, "/tasks-lists/", h.middlewares.ApplyMiddlewares(
		h.DeleteTasksList,
		h.middlewares.ForAuth,
	))
}

func (h TasksListsHandler) GetAllTasksLists(w http.ResponseWriter, r *http.Request) {
	user, err := h.Parent.getUser(r)
	if err != nil {
		h.Parent.error(w, fmt.Sprintf("can not get user: %s", err.Error()), http.StatusInternalServerError)
	}

	tasks, err := h.Services.TasksLists.GetAllUserTasksLists(context.Background(), user.ID)
	if err != nil {
		h.Parent.error(w, fmt.Sprintf("can not find user tasks: %s", err.Error()), http.StatusInternalServerError)
	}

	tasksListBytes, err := json.Marshal(tasks)
	if err != nil {
		h.Parent.error(w, fmt.Sprintf("can not find user tasks: %s", err.Error()), http.StatusInternalServerError)
	}

	h.Parent.send(w, string(tasksListBytes), http.StatusOK)
}

func (h TasksListsHandler) AddNewTasksList(w http.ResponseWriter, r *http.Request) {
	var CreateTasksListRB models.CreateTasksListRB
	var unmarshalErr *json.UnmarshalTypeError

	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()
	err := decoder.Decode(&CreateTasksListRB)
	if err != nil {
		if errors.As(err, &unmarshalErr) {
			h.Parent.error(w, fmt.Sprintf("bad Request: wrong type provided for field - %s", unmarshalErr.Field), http.StatusBadRequest)
		} else {
			h.Parent.error(w, fmt.Sprintf("bad Request: %s", err.Error()), http.StatusBadRequest)
		}
		return
	}

	if err := validator.Validate(CreateTasksListRB); err != nil {
		h.Parent.error(w, fmt.Sprintf("validataion error: %s", err.Error()), http.StatusBadRequest)
		return
	}

	user, err := h.Parent.getUser(r)
	if err != nil {
		h.Parent.error(w, fmt.Sprintf("can not find user: %s", err.Error()), http.StatusInternalServerError)
		return
	}

	CreateTasksListDTO := CreateTasksListRB.Build(user.ID)

	tasksList, err := h.Services.TasksLists.AddTasksList(context.Background(), CreateTasksListDTO)
	if err != nil {
		h.Parent.error(w, fmt.Sprintf("can not add tasks list: %s", err.Error()), http.StatusInternalServerError)
		return
	}

	taskBytes, _ := json.Marshal(tasksList)

	h.Parent.send(w, string(taskBytes), http.StatusOK)
}

func (h TasksListsHandler) UpdateTasksList(w http.ResponseWriter, r *http.Request) {
	var UpdateTasksListRB models.UpdateTasksListRB
	var unmarshalErr *json.UnmarshalTypeError

	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()
	err := decoder.Decode(&UpdateTasksListRB)
	if err != nil {
		if errors.As(err, &unmarshalErr) {
			h.Parent.error(w, fmt.Sprintf("bad Request: wrong type provided for field - %s", unmarshalErr.Field), http.StatusBadRequest)
		} else {
			h.Parent.error(w, fmt.Sprintf("bad Request: %s", err.Error()), http.StatusBadRequest)
		}
		return
	}

	if err := validator.Validate(UpdateTasksListRB); err != nil {
		h.Parent.error(w, fmt.Sprintf("validataion error: %s", err.Error()), http.StatusBadRequest)
		return
	}

	user, err := h.Parent.getUser(r)
	if err != nil {
		h.Parent.error(w, fmt.Sprintf("can not find user: %s", err.Error()), http.StatusInternalServerError)
		return
	}

	UpdateTasksListDTO := UpdateTasksListRB.Build()

	tasksList, err := h.Services.TasksLists.UpdateTasksList(context.Background(), UpdateTasksListRB.ID, user.ID, UpdateTasksListDTO)
	if err != nil {
		h.Parent.error(w, fmt.Sprintf("can not update tasks list: %s", err.Error()), http.StatusInternalServerError)
		return
	}

	tasksListBytes, _ := json.Marshal(tasksList)

	h.Parent.send(w, string(tasksListBytes), http.StatusOK)
}

func (h TasksListsHandler) DeleteTasksList(w http.ResponseWriter, r *http.Request) {
	var DeleteTasksListDTO models.DeleteTasksListDTO
	var unmarshalErr *json.UnmarshalTypeError

	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()
	err := decoder.Decode(&DeleteTasksListDTO)
	if err != nil {
		if errors.As(err, &unmarshalErr) {
			h.Parent.error(w, fmt.Sprintf("bad Request: wrong type provided for field - %s", unmarshalErr.Field), http.StatusBadRequest)
		} else {
			h.Parent.error(w, fmt.Sprintf("bad Request: %s", err.Error()), http.StatusBadRequest)
		}
		return
	}

	if err := validator.Validate(DeleteTasksListDTO); err != nil {
		h.Parent.error(w, fmt.Sprintf("validataion error: %s", err.Error()), http.StatusBadRequest)
		return
	}

	user, err := h.Parent.getUser(r)
	if err != nil {
		h.Parent.error(w, fmt.Sprintf("can not find user: %s", err.Error()), http.StatusInternalServerError)
		return
	}

	tlid, err := h.Services.TasksLists.DeleteTasksList(context.Background(), DeleteTasksListDTO.ID, user.ID)
	if err != nil {
		h.Parent.error(w, fmt.Sprintf("can not delete tasks list: %s", err.Error()), http.StatusInternalServerError)
		return
	}

	h.Parent.send(w, fmt.Sprintf("\"%s\"", tlid), http.StatusOK)
}
