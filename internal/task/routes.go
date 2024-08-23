package task

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/go-chi/chi/v5"
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
		r.Get("/{taskID}", h.handleGetTask)
		r.Put("/", h.handleUpdateTask)
		r.Delete("/{taskID}", h.handleDeleteTask)
	})
	router.Get("/api/task-containers/{containerID}/tasks", h.handleGetTasksByContainerId)
	router.Post("/api/task-containers/{containerID}/tasks", h.handleCreateTask)
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
	_, err = h.taskRepo.CreateTask(container.ContainerId, *createDto)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}
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
