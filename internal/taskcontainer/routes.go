package taskcontainer

import (
	"fmt"
	"net/http"

	"example.com/taskapp/internal/auth"
	"example.com/taskapp/internal/user"
	"example.com/taskapp/utils"
	"github.com/go-chi/chi/v5"
)

type Handler struct {
	containerRepo ContainerRepository
	userRepo      user.UserRepository
}

func NewHandler(repo ContainerRepository, userRepo user.UserRepository) *Handler {
	return &Handler{containerRepo: repo, userRepo: userRepo}
}
func (h *Handler) RegisterRoutes(router *chi.Mux) {
	router.Route("/api/task-containers", func(r chi.Router) {
		r.Get("/", h.handleGetTaskContainers)
		r.Get("/{containerID}", auth.WithJWTAuth(h.handleGetTaskContainerById))
	})
}
func (h *Handler) handleGetTaskContainers(w http.ResponseWriter, r *http.Request) {
	users, err := h.containerRepo.AllTaskContainers()
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}
	utils.WriteJSON(w, http.StatusOK, users)
}
func (h *Handler) handleGetTaskContainerById(w http.ResponseWriter, r *http.Request) {
	containerId := chi.URLParam(r, "containerID")
	if containerId == "" {
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("missing Container ID"))
		return
	}
	// TODO Requirement - check if container ID is UUID format.

	container, err := h.containerRepo.GetById(containerId)
	if err != nil {
		utils.WriteError(w, http.StatusNotFound, fmt.Errorf("container does not exist"))
		return
	}
	utils.WriteJSON(w, http.StatusOK, container)
}
