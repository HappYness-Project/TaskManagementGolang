package taskcontainer

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/happYness-Project/taskManagementGolang/internal/auth"
	"github.com/happYness-Project/taskManagementGolang/internal/user"
	"github.com/happYness-Project/taskManagementGolang/utils"
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
		// TODO Jwt Middleware setup in here, maybe more higher level if possible.
		r.Get("/", auth.WithJWTAuth(h.handleGetTaskContainers))
		r.Get("/{containerID}", auth.WithJWTAuth(h.handleGetTaskContainerById))
	})
	router.Get("/api/user-groups/{usergroupID}/task-containers", auth.WithJWTAuth(h.handleGetTaskContainersByGroupId))
}
func (h *Handler) handleGetTaskContainers(w http.ResponseWriter, r *http.Request) {
	users, err := h.containerRepo.AllTaskContainers()
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}
	utils.WriteJsonWithEncode(w, http.StatusOK, users)
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
	utils.WriteJsonWithEncode(w, http.StatusOK, container)
}

func (h *Handler) handleGetTaskContainersByGroupId(w http.ResponseWriter, r *http.Request) {
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

	containers, err := h.containerRepo.GetContainersByGroupId(groupId)
	if err != nil {
		utils.WriteError(w, http.StatusNotFound, fmt.Errorf("container does not exist"))
		return
	}
	containersJson, _ := json.Marshal(containers)
	utils.WriteJSON(w, http.StatusOK, containersJson)
}
