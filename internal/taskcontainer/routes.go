package taskcontainer

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
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
		r.Post("/", auth.WithJWTAuth(h.handleCreateTaskContainer))
		r.Get("/", auth.WithJWTAuth(h.handleGetTaskContainers))
		r.Get("/{containerID}", auth.WithJWTAuth(h.handleGetTaskContainerById))
		r.Delete("/{containerID}", auth.WithJWTAuth(h.handleDeleteTaskContainer))
	})
	router.Get("/api/user-groups/{usergroupID}/task-containers", auth.WithJWTAuth(h.handleGetTaskContainersByGroupId))
}
func (h *Handler) handleGetTaskContainers(w http.ResponseWriter, r *http.Request) {
	containers, err := h.containerRepo.AllTaskContainers()
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}
	containersJson, _ := json.Marshal(containers)
	utils.WriteJsonWithEncode(w, http.StatusOK, containersJson)
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
	utils.WriteJsonWithEncode(w, http.StatusOK, containersJson)
}

func (h *Handler) handleCreateTaskContainer(w http.ResponseWriter, r *http.Request) {
	var createDto CreateContainerDto
	if err := utils.ParseJSON(r, &createDto); err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}

	container := TaskContainer{
		Id:             uuid.New().String(),
		Name:           createDto.Name,
		Description:    createDto.Description,
		Type:           createDto.Type,
		IsActive:       true,
		activity_level: 0,
		UsergroupId:    createDto.UserGroupId,
	}
	_ = h.containerRepo.CreateContainer(container)
	utils.WriteJsonWithEncode(w, http.StatusCreated, container.Id)
}

func (h *Handler) handleDeleteTaskContainer(w http.ResponseWriter, r *http.Request) {
	containerId := chi.URLParam(r, "containerID")
	err := h.containerRepo.DeleteContainer(containerId)
	if err != nil {
		utils.WriteError(w, http.StatusNotFound, fmt.Errorf("container does not exist"))
		return
	}
	utils.WriteJsonWithEncode(w, http.StatusNoContent, "task container is removed.")
}
