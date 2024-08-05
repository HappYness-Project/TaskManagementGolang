package task

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"example.com/taskapp/utils"
	"github.com/go-chi/chi/v5"
)

type Handler struct {
	taskRepo TaskRepository
}

func NewHandler(repo TaskRepository) *Handler {
	return &Handler{taskRepo: repo}
}
func (h *Handler) RegisterRoutes(router *chi.Mux) {
	router.Route("/api/tasks", func(r chi.Router) {
		r.Get("/", h.handleGetTasks)
		r.Get("/{taskID}", h.handleGetTask)
		r.Put("/", h.handleUpdateTask)
		r.Delete("/", h.handleDeleteTask)
	})
	router.Get("/api/task-containers/{containerID}/tasks", h.handleGetTaskByContainerId)
	router.Post("/api/task-containers/{containerID}/tasks", h.handleCreateTask)
}
func (h *Handler) handleGetTasks(w http.ResponseWriter, r *http.Request) {
	tasks, err := h.taskRepo.GetAllTasks()
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}
	utils.WriteJSON(w, http.StatusOK, tasks)
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
	utils.WriteJSON(w, http.StatusOK, task)
}

func (h *Handler) handleGetTaskByContainerId(w http.ResponseWriter, r *http.Request) {
	// tasks, err := h.taskRepo.GetAllTasksByContainerId()
}

func (h *Handler) handleCreateTask(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Error reading request body", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	var task *Task
	err = json.Unmarshal(body, &task)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
	}
}
func (h *Handler) handleUpdateTask(w http.ResponseWriter, r *http.Request) {
}

func (h *Handler) handleDeleteTask(w http.ResponseWriter, r *http.Request) {
}
