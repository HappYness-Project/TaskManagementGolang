package task

import (
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
}
func (h *Handler) handleGetTasks(w http.ResponseWriter, r *http.Request) {
	tasks, err := h.taskRepo.GetAllTasks()
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}
	utils.WriteJSON(w, http.StatusOK, tasks)
}
