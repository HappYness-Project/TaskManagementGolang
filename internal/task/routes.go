package task

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/happYness-Project/taskManagementGolang/internal/auth"
	"github.com/happYness-Project/taskManagementGolang/internal/taskcontainer"
	"github.com/happYness-Project/taskManagementGolang/utils"
)

type Handler struct {
	taskRepo      TaskRepository
	containerRepo taskcontainer.ContainerRepository
}

func NewHandler(repo TaskRepository, tcRepo taskcontainer.ContainerRepository) *Handler {
	return &Handler{taskRepo: repo, containerRepo: tcRepo}
}
func (h *Handler) RegisterRoutes(router *chi.Mux) {
	router.Route("/api/tasks", func(r chi.Router) {
		r.Get("/", h.handleGetTasks)
		r.Get("/{taskID}", auth.WithJWTAuth(h.handleGetTask))
		r.Put("/", auth.WithJWTAuth(h.handleUpdateTask))
		r.Delete("/{taskID}", auth.WithJWTAuth(h.handleDeleteTask))
		r.Patch("/{taskID}/toggle-completion", auth.WithJWTAuth(h.handleDoneTask))
	})
	router.Get("/api/task-containers/{containerID}/tasks", auth.WithJWTAuth(h.handleGetTasksByContainerId))
	router.Post("/api/task-containers/{containerID}/tasks", auth.WithJWTAuth(h.handleCreateTask))
}
func (h *Handler) handleGetTasks(w http.ResponseWriter, r *http.Request) {
	tasks, err := h.taskRepo.GetAllTasks()
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}
	utils.WriteJsonWithEncode(w, http.StatusOK, tasks)
}
func (h *Handler) handleGetTask(w http.ResponseWriter, r *http.Request) {
	taskId := chi.URLParam(r, "taskID")
	if taskId == "" {
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("missing Task ID"))
		return
	}
	// TODO Requirement - check if container ID is UUID format.
	task, err := h.taskRepo.GetTaskById(taskId)
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
	containerId := chi.URLParam(r, "containerID")
	if containerId == "" {
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("missing container ID"))
		return
	}
	container, err := h.containerRepo.GetById(containerId)
	if err != nil {
		utils.WriteError(w, http.StatusNotFound, fmt.Errorf("not able to find container"))
		return
	}
	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Error reading request body", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	var createDto *CreateTaskDto
	err = json.Unmarshal(body, &createDto)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}
	uuid, err := h.taskRepo.CreateTask(container.ContainerId, *createDto)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}
	utils.WriteJsonWithEncode(w, http.StatusCreated, uuid)
}

func (h *Handler) handleUpdateTask(w http.ResponseWriter, r *http.Request) {
}

func (h *Handler) handleDeleteTask(w http.ResponseWriter, r *http.Request) {
	taskId := chi.URLParam(r, "taskID")
	if taskId == "" {
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("missing Task ID"))
		return
	}
	err := h.taskRepo.DeleteTask(taskId)
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
