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

type TasksHandler struct {
	Parent      *Router
	Router      *httprouter.Router
	Services    *services.Services
	middlewares *middlewares.Middlewares
	logger      *logging.Logger
	redis       *redis.Client
}

func NewTasksHandler(router *Router) *TasksHandler {
	return &TasksHandler{
		Parent:      router,
		Router:      router.Router,
		Services:    router.Services,
		middlewares: router.middlewares,
		logger:      router.logger,
		redis:       router.redis,
	}
}

func (h TasksHandler) RegisterTasksRoutes() {
	h.Router.HandlerFunc(http.MethodGet, "/tasks/", h.middlewares.ApplyMiddlewares(
		h.GetAllTasks,
		h.middlewares.ForAuth,
	))
	h.Router.HandlerFunc(http.MethodPost, "/tasks/", h.middlewares.ApplyMiddlewares(
		h.AddNewTask,
		h.middlewares.ForAuth,
	))
	h.Router.HandlerFunc(http.MethodPatch, "/tasks/", h.middlewares.ApplyMiddlewares(
		h.UpdateTask,
		h.middlewares.ForAuth,
	))
	h.Router.HandlerFunc(http.MethodDelete, "/tasks/", h.middlewares.ApplyMiddlewares(
		h.DeleteTask,
		h.middlewares.ForAuth,
	))
	h.Router.HandlerFunc(http.MethodDelete, "/tasks/clear", h.middlewares.ApplyMiddlewares(
		h.DeleteAllTask,
		h.middlewares.ForAuth,
	))
}

func (h TasksHandler) GetAllTasks(w http.ResponseWriter, r *http.Request) {
	user, err := h.Parent.getUser(r)
	if err != nil {
		h.Parent.error(w, fmt.Sprintf("can not get user: %s", err.Error()), http.StatusInternalServerError)
	}

	tasks, err := h.Services.Tasks.GetAllUserTasks(context.Background(), user.ID)
	if err != nil {
		h.Parent.error(w, fmt.Sprintf("can not find user tasks: %s", err.Error()), http.StatusInternalServerError)
	}

	tasksBytes, err := json.Marshal(tasks)
	if err != nil {
		h.Parent.error(w, fmt.Sprintf("can not find user tasks: %s", err.Error()), http.StatusInternalServerError)
	}

	h.Parent.send(w, string(tasksBytes), http.StatusOK)
}

func (h TasksHandler) AddNewTask(w http.ResponseWriter, r *http.Request) {

	var createTaskDTO models.CreateTaskDTO
	var unmarshalErr *json.UnmarshalTypeError

	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()
	err := decoder.Decode(&createTaskDTO)
	if err != nil {
		if errors.As(err, &unmarshalErr) {
			h.Parent.error(w, fmt.Sprintf("bad Request: wrong type provided for field - %s", unmarshalErr.Field), http.StatusBadRequest)
		} else {
			h.Parent.error(w, fmt.Sprintf("bad Request: %s", err.Error()), http.StatusBadRequest)
		}
		return
	}

	if err := validator.Validate(createTaskDTO); err != nil {
		h.Parent.error(w, fmt.Sprintf("validataion error: %s", err.Error()), http.StatusBadRequest)
		return
	}

	user, err := h.Parent.getUser(r)
	if err != nil {
		h.Parent.error(w, fmt.Sprintf("can not find user: %s", err.Error()), http.StatusInternalServerError)
		return
	}

	createTaskDTO.UserID = user.ID

	task, err := h.Services.Tasks.AddTask(context.Background(), createTaskDTO)
	if err != nil {
		h.Parent.error(w, fmt.Sprintf("can not add task: %s", err.Error()), http.StatusInternalServerError)
		return
	}

	taskBytes, _ := json.Marshal(task)

	h.Parent.send(w, string(taskBytes), http.StatusOK)
}

func (h TasksHandler) UpdateTask(w http.ResponseWriter, r *http.Request) {
	var taskDTO models.Task
	var unmarshalErr *json.UnmarshalTypeError

	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()
	err := decoder.Decode(&taskDTO)
	if err != nil {
		if errors.As(err, &unmarshalErr) {
			h.Parent.error(w, fmt.Sprintf("bad Request: wrong type provided for field - %s", unmarshalErr.Field), http.StatusBadRequest)
		} else {
			h.Parent.error(w, fmt.Sprintf("bad Request: %s", err.Error()), http.StatusBadRequest)
		}
		return
	}

	updateTaskDTO := taskDTO.BuildForUpdate()

	if err := validator.Validate(updateTaskDTO); err != nil {
		h.Parent.error(w, fmt.Sprintf("validataion error: %s", err.Error()), http.StatusBadRequest)
		return
	}

	user, err := h.Parent.getUser(r)
	if err != nil {
		h.Parent.error(w, fmt.Sprintf("can not find user: %s", err.Error()), http.StatusInternalServerError)
		return
	}

	task, err := h.Services.Tasks.UpdateTask(context.Background(), taskDTO.ID, user.ID, updateTaskDTO)
	if err != nil {
		h.Parent.error(w, fmt.Sprintf("can not update task: %s", err.Error()), http.StatusInternalServerError)
		return
	}

	taskBytes, _ := json.Marshal(task)

	h.Parent.send(w, string(taskBytes), http.StatusOK)
}

func (h TasksHandler) DeleteTask(w http.ResponseWriter, r *http.Request) {
	var deleteTaskDTO models.DeleteTaskDTO
	var unmarshalErr *json.UnmarshalTypeError

	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()
	err := decoder.Decode(&deleteTaskDTO)
	if err != nil {
		if errors.As(err, &unmarshalErr) {
			h.Parent.error(w, fmt.Sprintf("bad Request: wrong type provided for field - %s", unmarshalErr.Field), http.StatusBadRequest)
		} else {
			h.Parent.error(w, fmt.Sprintf("bad Request: %s", err.Error()), http.StatusBadRequest)
		}
		return
	}

	if err := validator.Validate(deleteTaskDTO); err != nil {
		h.Parent.error(w, fmt.Sprintf("validataion error: %s", err.Error()), http.StatusBadRequest)
		return
	}

	user, err := h.Parent.getUser(r)
	if err != nil {
		h.Parent.error(w, fmt.Sprintf("can not find user: %s", err.Error()), http.StatusInternalServerError)
		return
	}

	tid, err := h.Services.Tasks.DeleteTask(context.Background(), deleteTaskDTO.ID, user.ID)
	if err != nil {
		h.Parent.error(w, fmt.Sprintf("can not delete task: %s", err.Error()), http.StatusInternalServerError)
		return
	}

	h.Parent.send(w, fmt.Sprintf("\"%s\"", tid), http.StatusOK)
}

func (h TasksHandler) DeleteAllTask(w http.ResponseWriter, r *http.Request) {
	user, err := h.Parent.getUser(r)
	if err != nil {
		h.Parent.error(w, fmt.Sprintf("can not find user: %s", err.Error()), http.StatusInternalServerError)
		return
	}

	err = h.Services.Tasks.DeleteAllTask(context.Background(), user.ID)
	if err != nil {
		h.Parent.error(w, fmt.Sprintf("can not delete task: %s", err.Error()), http.StatusInternalServerError)
		return
	}

	h.Parent.send(w, "\"\"", http.StatusOK)
}
