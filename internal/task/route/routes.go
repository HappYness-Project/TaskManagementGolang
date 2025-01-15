package route

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/happYness-Project/taskManagementGolang/internal/task/model"
	taskRepo "github.com/happYness-Project/taskManagementGolang/internal/task/repository"
	containerRepo "github.com/happYness-Project/taskManagementGolang/internal/taskcontainer/repository"
	"github.com/happYness-Project/taskManagementGolang/internal/usergroup"
	"github.com/happYness-Project/taskManagementGolang/pkg/loggers"
	"github.com/happYness-Project/taskManagementGolang/pkg/utils"
)

type Handler struct {
	logger        *loggers.AppLogger
	taskRepo      taskRepo.TaskRepository
	containerRepo containerRepo.ContainerRepository
	groupRepo     usergroup.UserGroupRepository
}

func NewHandler(logger *loggers.AppLogger, repo taskRepo.TaskRepository, tcRepo containerRepo.ContainerRepository, ugRepo usergroup.UserGroupRepository) *Handler {
	return &Handler{logger: logger, taskRepo: repo, containerRepo: tcRepo, groupRepo: ugRepo}
}
func (h *Handler) RegisterRoutes(router *chi.Mux) {
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
		h.logger.Error().
			Err(err).
			Str("ErrorCode", TaskGetServerError).
			Msg("Error occurred during GetAllTasks.")
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}
	utils.WriteJsonWithEncode(w, http.StatusOK, tasks)
}
func (h *Handler) handleGetTask(w http.ResponseWriter, r *http.Request) {
	// TODO Requirement - check if container ID is UUID format.
	task, err := h.taskRepo.GetTaskById(chi.URLParam(r, "taskID"))
	if err != nil {
		utils.WriteError(w, http.StatusNotFound, fmt.Errorf("task does not exist"))
		return
	}
	utils.WriteJsonWithEncode(w, http.StatusOK, task)
}

func (h *Handler) handleGetTasksByContainerId(w http.ResponseWriter, r *http.Request) {
	containerId := chi.URLParam(r, "containerID")
	if containerId == "" {
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("missing container ID"))
		return
	}
	tasks, err := h.taskRepo.GetTasksByContainerId(containerId)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("getting tasks failure"))
		return
	}
	utils.WriteJsonWithEncode(w, http.StatusOK, tasks)
}

func (h *Handler) handleCreateTask(w http.ResponseWriter, r *http.Request) {
	container, err := h.containerRepo.GetById(chi.URLParam(r, "containerID"))
	if err != nil {
		utils.WriteError(w, http.StatusNotFound, fmt.Errorf("not able to find container"))
		return
	}
	//TODO Validation logic
	var createDto CreateTaskDto
	if err := utils.ParseJson(r, &createDto); err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
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
	uuid, err := h.taskRepo.CreateTask(container.Id, task)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}
	utils.WriteJsonWithEncode(w, http.StatusCreated, uuid)
}

func (h *Handler) handleUpdateTask(w http.ResponseWriter, r *http.Request) {
	taskId := chi.URLParam(r, "taskID")
	var updateDto UpdateTaskDto
	if err := utils.ParseJson(r, &updateDto); err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}
	task, err := h.taskRepo.GetTaskById(taskId)
	if err != nil || task == nil {
		utils.WriteError(w, http.StatusNotFound, fmt.Errorf("cannot find task"))
		return
	}

	task.TaskName = updateDto.TaskName
	task.TaskDesc = updateDto.TaskDesc
	task.TargetDate = updateDto.TargetDate
	task.Priority = updateDto.Priority
	task.Category = updateDto.Category
	err = h.taskRepo.UpdateTask(*task)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("not able to update task"))
		return
	}

}

func (h *Handler) handleDeleteTask(w http.ResponseWriter, r *http.Request) {
	err := h.taskRepo.DeleteTask(chi.URLParam(r, "taskID"))
	if err != nil {
		utils.WriteError(w, http.StatusNotFound, fmt.Errorf("error occurred during deleting a task"))
		return
	}
	utils.WriteJsonWithEncode(w, http.StatusNoContent, "task has been removed.")
}

func (h *Handler) handleDoneTask(w http.ResponseWriter, r *http.Request) {
	taskId := chi.URLParam(r, "taskID")
	if taskId == "" {
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("missing Task ID"))
		return
	}
	task, _ := h.taskRepo.GetTaskById(taskId)
	if task == nil {
		utils.WriteError(w, http.StatusNotFound, fmt.Errorf("not found task"))
		return
	}

	type ToggleBody struct {
		IsCompleted bool `json:"is_completed"`
	}
	var toggleBody ToggleBody
	err := json.NewDecoder(r.Body).Decode(&toggleBody)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = h.taskRepo.DoneTask(taskId, toggleBody.IsCompleted)
	if err != nil {
		utils.WriteError(w, http.StatusNotFound, fmt.Errorf("error occurred during done task"))
		return
	}
	utils.WriteJsonWithEncode(w, http.StatusOK, "task is changed to Done.")
}
func (h *Handler) handleImportantTask(w http.ResponseWriter, r *http.Request) {
	taskId := chi.URLParam(r, "taskID")
	if taskId == "" {
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("missing Task ID"))
		return
	}
	task, _ := h.taskRepo.GetTaskById(taskId)
	if task == nil {
		utils.WriteError(w, http.StatusNotFound, fmt.Errorf("not found task"))
		return
	}

	type ToggleBody struct {
		IsImportant bool `json:"is_important"`
	}
	var toggleBody ToggleBody
	err := json.NewDecoder(r.Body).Decode(&toggleBody)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = h.taskRepo.UpdateImportantTask(taskId, toggleBody.IsImportant)
	if err != nil {
		utils.WriteError(w, http.StatusNotFound, fmt.Errorf("error occurred during important toggle task"))
		return
	}
	utils.WriteJsonWithEncode(w, http.StatusOK, "task important field is changed.")
}

func (h *Handler) handleGetTasksByGroupId(w http.ResponseWriter, r *http.Request) {
	groupIdVar := chi.URLParam(r, "usergroupID")
	if groupIdVar == "" {
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("missing Group ID"))
		return
	}
	groupId, err := strconv.Atoi(groupIdVar)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("invalid Group ID"))
		return
	}
	usergroup, err := h.groupRepo.GetById(groupId)
	if err != nil {
		utils.WriteError(w, http.StatusNotFound, fmt.Errorf("usergroup cannot be found"))
		return
	}
	var tasks []model.Task

	if r.URL.Query().Get("important") == "true" {
		tasks, err = h.taskRepo.GetAllTasksByGroupIdOnlyImportant(usergroup.GroupId)
		if err != nil {
			utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("error ocurred during getting important tasks"))
			return
		}
	} else if r.URL.Query().Get("important") == "false" {
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("error ocurred during getting tasks - not implemented for false case"))
		return
	} else {
		tasks, err = h.taskRepo.GetAllTasksByGroupId(usergroup.GroupId)
		if err != nil {
			utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("error ocurred during getting tasks"))
			return
		}
	}
	utils.WriteJsonWithEncode(w, http.StatusOK, tasks)
}
