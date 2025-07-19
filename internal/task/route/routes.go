package route

import (
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/happYness-Project/taskManagementGolang/internal/task/model"
	taskRepo "github.com/happYness-Project/taskManagementGolang/internal/task/repository"
	containerRepo "github.com/happYness-Project/taskManagementGolang/internal/taskcontainer/repository"
	usergroupRepo "github.com/happYness-Project/taskManagementGolang/internal/usergroup/repository"
	usergroupRoute "github.com/happYness-Project/taskManagementGolang/internal/usergroup/route"
	"github.com/happYness-Project/taskManagementGolang/pkg/constants"
	"github.com/happYness-Project/taskManagementGolang/pkg/loggers"
	"github.com/happYness-Project/taskManagementGolang/pkg/response"
)

type Handler struct {
	logger        *loggers.AppLogger
	taskRepo      taskRepo.TaskRepository
	containerRepo containerRepo.ContainerRepository
	groupRepo     usergroupRepo.UserGroupRepository
}

func NewHandler(logger *loggers.AppLogger, repo taskRepo.TaskRepository, tcRepo containerRepo.ContainerRepository, ugRepo usergroupRepo.UserGroupRepository) *Handler {
	return &Handler{logger: logger, taskRepo: repo, containerRepo: tcRepo, groupRepo: ugRepo}
}
func (h *Handler) RegisterRoutes(router chi.Router) {
	router.Route("/api/tasks", func(r chi.Router) {
		r.Get("/", h.handleGetTasks)
		r.Get("/{taskID}", h.handleGetTask)
		r.Put("/{taskID}", h.handleUpdateTask)
		r.Delete("/{taskID}", h.handleDeleteTask)
		r.Patch("/{taskID}/toggle-completion", h.handleDoneTask)
		r.Patch("/{taskID}/toggle-important", h.handleImportantTask)
	})
	router.Get("/api/task-containers/{containerID}/tasks", h.handleGetTasksByContainerId)
	router.Post("/api/task-containers/{containerID}/tasks", h.handleCreateTask)
	router.Get("/api/user-groups/{usergroupID}/tasks", h.handleGetTasksByGroupId)
}
func (h *Handler) handleGetTasks(w http.ResponseWriter, r *http.Request) {
	tasks, err := h.taskRepo.GetAllTasks()
	if err != nil {
		h.logger.Error().Err(err).Str("ErrorCode", TaskGetServerError).Msg(err.Error())
		response.InternalServerError(w, "Error occurred during getting all tasks.")
		return
	}
	response.SuccessJson(w, tasks, "successfully get tasks", http.StatusOK)
}
func (h *Handler) handleGetTask(w http.ResponseWriter, r *http.Request) {
	task, err := h.taskRepo.GetTaskById(chi.URLParam(r, "taskID"))
	if err != nil {
		h.logger.Error().Err(err).Str("ErrorCode", TaskGetNotFound).Msg("Error occurred during GetTask.")
		response.ErrorResponse(w, http.StatusNotFound, *(response.New(TaskGetNotFound, "Not Found", "task does not exist")))
		return
	}
	response.WriteJsonWithEncode(w, http.StatusOK, task)
}

func (h *Handler) handleGetTasksByContainerId(w http.ResponseWriter, r *http.Request) {
	containerId := chi.URLParam(r, "containerID")
	if containerId == "" {
		h.logger.Error().Msg("container Id missing")
		response.BadRequestMissingParameters(w)
		return
	}
	tasks, err := h.taskRepo.GetTasksByContainerId(containerId)
	if err != nil {
		h.logger.Error().Err(err).Str("ErrorCode", TaskGetServerError).Msg("Error occurred during GetTasksByContainerId")
		response.ErrorResponse(w, http.StatusInternalServerError, *(response.New(TaskGetServerError, "Failed to get tasks by container id", err.Error())))
		return
	}
	response.WriteJsonWithEncode(w, http.StatusOK, tasks)
}

func (h *Handler) handleCreateTask(w http.ResponseWriter, r *http.Request) {
	containerId := chi.URLParam(r, "containerID")
	if containerId == "" {
		h.logger.Error().Msg("container Id missing")
		response.BadRequestMissingParameters(w)
		return
	}
	var createDto CreateTaskDto
	if err := response.ParseJson(r, &createDto); err != nil {
		h.logger.Error().Err(err).Str("ErrorCode", constants.RequestBodyError).Msg("Invalid JSON body for CreateTaskDto")
		response.ErrorResponse(w, http.StatusBadRequest, *(response.New(constants.RequestBodyError, "Invalid Json Body", err.Error())))
		return
	}
	container, err := h.containerRepo.GetById(containerId)
	if err != nil {
		h.logger.Error().Err(err).Str("ErrorCode", TaskGetTaskContainerNotFound).Msg("Task container not found")
		response.ErrorResponse(w, http.StatusNotFound, *(response.New(TaskGetTaskContainerNotFound, "Task container not found", err.Error())))
		return
	}

	task := model.Task{
		TaskId:     uuid.New().String(),
		TaskName:   createDto.TaskName,
		TaskDesc:   createDto.TaskDesc,
		TaskType:   "",
		CreatedAt:  time.Now(),
		UpdatedAt:  time.Now(),
		TargetDate: createDto.TargetDate,
		Priority:   createDto.Priority,
		Category:   createDto.Category,
	}
	newTask, err := h.taskRepo.CreateTask(container.Id, task)
	if err != nil {
		h.logger.Error().Err(err).Str("ErrorCode", TaskCreateServerError).Msg("Error occurred during CreateTask")
		response.ErrorResponse(w, http.StatusBadRequest, *(response.New(TaskCreateServerError, "Failed to create task", err.Error())))
		return
	}
	response.WriteJsonWithEncode(w, http.StatusCreated, newTask)
}

func (h *Handler) handleUpdateTask(w http.ResponseWriter, r *http.Request) {
	var updateDto UpdateTaskDto
	if err := response.ParseJson(r, &updateDto); err != nil {
		h.logger.Error().Err(err).Str("ErrorCode", constants.RequestBodyError).Msg("Invalid JSON body for UpdateTaskDto")
		response.ErrorResponse(w, http.StatusBadRequest, *(response.New(constants.RequestBodyError, "Request Body Error", err.Error())))
		return
	}
	taskId := chi.URLParam(r, "taskID")
	task, err := h.taskRepo.GetTaskById(taskId)
	if err != nil || task == nil {
		h.logger.Error().Err(err).Str("ErrorCode", TaskGetNotFound).Msg("Cannot find task for update")
		response.NotFound(w, TaskGetNotFound, "Cannot find task")
		return
	}

	task.TaskName = updateDto.TaskName
	task.TaskDesc = updateDto.TaskDesc
	task.TargetDate = updateDto.TargetDate
	task.Priority = updateDto.Priority
	task.Category = updateDto.Category
	err = h.taskRepo.UpdateTask(*task)
	if err != nil {
		h.logger.Error().Err(err).Str("ErrorCode", TaskUpdateServerError).Msg("Not able to update task")
		response.ErrorResponse(w, http.StatusBadRequest, *(response.New(TaskUpdateServerError, "Failed to update task", err.Error())))
		return
	}
	response.WriteJsonWithEncode(w, http.StatusOK, "task updated successfully")
}

func (h *Handler) handleDeleteTask(w http.ResponseWriter, r *http.Request) {
	taskId := chi.URLParam(r, "taskID")
	err := h.taskRepo.DeleteTask(taskId)
	if err != nil {
		// TODO : Different types of error. need to identify
		h.logger.Error().Err(err).Str("ErrorCode", TaskDeleteServerError).Msg("Error occurred during deleting a task")
		response.ErrorResponse(w, http.StatusNotFound, *(response.New(TaskDeleteServerError, "Failed to delete task", err.Error())))
		return
	}
	response.WriteJsonWithEncode(w, http.StatusNoContent, "task has been removed.")
}

func (h *Handler) handleDoneTask(w http.ResponseWriter, r *http.Request) {
	taskId := chi.URLParam(r, "taskID")
	if taskId == "" {
		h.logger.Error().Msg("missing Task ID")
		response.BadRequestMissingParameters(w) // needs to be fixed - Adding missing task id as details.
		return
	}
	task, _ := h.taskRepo.GetTaskById(taskId)
	if task == nil {
		h.logger.Error().Msg("not found task")
		response.NotFound(w, TaskGetNotFound, "Task not found")
		return
	}

	type ToggleBody struct {
		IsCompleted bool `json:"is_completed"`
	}
	var toggleBody ToggleBody
	err := json.NewDecoder(r.Body).Decode(&toggleBody)
	if err != nil {
		h.logger.Error().Err(err).Msg("Invalid JSON body for toggle completion")
		response.InvalidJsonBody(w, "Invalid Json body for toggle completion")
		return
	}

	err = h.taskRepo.DoneTask(taskId, toggleBody.IsCompleted)
	if err != nil {
		h.logger.Error().Err(err).Str("ErrorCode", TaskStatusDoneError).Msg("Error occurred during done task")
		response.ErrorResponse(w, http.StatusNotFound, *(response.New(TaskStatusDoneError, "Failed to toggle done")))
		return
	}
	response.WriteJsonWithEncode(w, http.StatusOK, "task is changed to Done.")
}

func (h *Handler) handleImportantTask(w http.ResponseWriter, r *http.Request) {
	taskId := chi.URLParam(r, "taskID")
	if taskId == "" {
		h.logger.Error().Msg("missing Task ID")
		response.BadRequestMissingParameters(w, "missing task id")
		return
	}
	task, _ := h.taskRepo.GetTaskById(taskId)
	if task == nil {
		h.logger.Error().Msg("not found task")
		response.NotFound(w, TaskGetNotFound, "task not found")
		return
	}

	type ToggleBody struct {
		IsImportant bool `json:"is_important"`
	}
	var toggleBody ToggleBody
	err := json.NewDecoder(r.Body).Decode(&toggleBody)
	if err != nil {
		h.logger.Error().Err(err).Msg("Invalid JSON body for toggle important")
		response.InvalidJsonBody(w, "Invalid json body for toggle important")
		return
	}

	err = h.taskRepo.UpdateImportantTask(taskId, toggleBody.IsImportant)
	if err != nil {
		h.logger.Error().Err(err).Str("ErrorCode", TaskUpdateImportantError).Msg("Error occurred during important toggle task")
		response.ErrorResponse(w, http.StatusBadRequest, *(response.New(TaskUpdateImportantError, "Failed to toggle important")))
		return
	}
	response.WriteJsonWithEncode(w, http.StatusOK, "task important field is changed.")
}

func (h *Handler) handleGetTasksByGroupId(w http.ResponseWriter, r *http.Request) {
	groupIdVar := chi.URLParam(r, "usergroupID")
	if groupIdVar == "" {
		h.logger.Error().Msg("missing Group ID")
		response.BadRequestMissingParameters(w, "missing group id")
		return
	}
	groupId, err := strconv.Atoi(groupIdVar)
	if err != nil {
		h.logger.Error().Err(err).Msg("invalid Group ID")
		response.ErrorResponse(w, http.StatusBadRequest, *(response.New(constants.InvalidParameter, "Invalid Group ID")))
		return
	}
	usergroup, err := h.groupRepo.GetById(groupId)
	if err != nil {
		h.logger.Error().Err(err).Msg("usergroup cannot be found")
		response.NotFound(w, usergroupRoute.UserGroupGetNotFound, "usergroup cannot be found")
		return
	}
	var tasks []model.Task

	if r.URL.Query().Get("important") == "true" {
		tasks, err = h.taskRepo.GetAllTasksByGroupIdOnlyImportant(usergroup.GroupId)
		if err != nil {
			h.logger.Error().Err(err).Msg("error occurred during getting important tasks")
			response.ErrorResponse(w, http.StatusBadRequest, *(response.New(TaskGetServerError, "Failed to get important tasks", err.Error())))
			return
		}
	} else if r.URL.Query().Get("important") == "false" {
		h.logger.Error().Msg("not implemented for false case")
		response.ErrorResponse(w, http.StatusBadRequest, *(response.New(TaskGetServerError, "Not implemented for important=false case")))
		return
	} else {
		tasks, err = h.taskRepo.GetAllTasksByGroupId(usergroup.GroupId)
		if err != nil {
			h.logger.Error().Err(err).Msg("error occurred during getting tasks")
			response.ErrorResponse(w, http.StatusBadRequest, *(response.New(TaskGetServerError, "Failed to get tasks")))
			return
		}
	}
	response.WriteJsonWithEncode(w, http.StatusOK, tasks)
}
