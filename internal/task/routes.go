package task

import (
	"encoding/json"
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
	router.Get("/tasks", h.handleGetTasks)
	// router.Get("/task-containers/{containerID}/tasks", h.handleGetTasksByContainerId)
}
func (h *Handler) handleGetTasks(w http.ResponseWriter, r *http.Request) {
	tasks, err := h.taskRepo.GetAllTasks()
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}
	utils.WriteJSON(w, http.StatusOK, tasks)
}

// func (h *Handler) handleGetTasks(w http.ResponseWriter, r *http.Request) {
// 	tasks, err := h.taskRepo.GetAllTasksByContainerId()
// }

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
