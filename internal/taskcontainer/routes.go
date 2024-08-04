package taskcontainer

import (
	"net/http"

	"example.com/taskapp/utils"
	"github.com/go-chi/chi/v5"
)

type Handler struct {
	containerRepo ContainerRepository
}

func NewHandler(repo ContainerRepository) *Handler {
	return &Handler{containerRepo: repo}
}
func (h *Handler) RegisterRoutes(router *chi.Mux) {
	router.Get("/task-containers", h.handleGetTaskContainers)
	// router.HandleFunc("/users/{userId}", h.handleGetUser).Methods(http.MethodGet)
}
func (h *Handler) handleGetTaskContainers(w http.ResponseWriter, r *http.Request) {
	users, err := h.containerRepo.AllTaskContainers()
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}
	utils.WriteJSON(w, http.StatusOK, users)
}
