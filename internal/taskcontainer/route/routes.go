package route

import (
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/happYness-Project/taskManagementGolang/internal/taskcontainer/model"
	container "github.com/happYness-Project/taskManagementGolang/internal/taskcontainer/repository"
	user "github.com/happYness-Project/taskManagementGolang/internal/user/repository"
	"github.com/happYness-Project/taskManagementGolang/pkg/constants"
	"github.com/happYness-Project/taskManagementGolang/pkg/loggers"
	"github.com/happYness-Project/taskManagementGolang/pkg/response"
)

type Handler struct {
	logger        *loggers.AppLogger
	containerRepo container.ContainerRepository
	userRepo      user.UserRepository
}

func NewHandler(logger *loggers.AppLogger, repo container.ContainerRepository, userRepo user.UserRepository) *Handler {
	return &Handler{logger: logger, containerRepo: repo, userRepo: userRepo}
}
func (h *Handler) RegisterRoutes(router chi.Router) {
	router.Route("/api/task-containers", func(r chi.Router) {
		r.Post("/", h.handleCreateTaskContainer)
		r.Get("/", h.handleGetTaskContainers)
		r.Get("/{containerID}", h.handleGetTaskContainerById)
		r.Delete("/{containerID}", h.handleDeleteTaskContainer)
	})
	router.Get("/api/user-groups/{usergroupID}/task-containers", h.handleGetTaskContainersByGroupId)
}
func (h *Handler) handleGetTaskContainers(w http.ResponseWriter, r *http.Request) {
	containers, err := h.containerRepo.AllTaskContainers()
	if err != nil {
		h.logger.Error().Err(err).Str("ErrorCode", TaskContainerGetError).Msg(err.Error())
		response.InternalServerError(w, "Error occurred during getting all task containers.")
		return
	}
	response.WriteJsonWithEncode(w, http.StatusOK, containers)
}
func (h *Handler) handleGetTaskContainerById(w http.ResponseWriter, r *http.Request) {
	containerId := chi.URLParam(r, "containerID")
	if containerId == "" {
		h.logger.Error().Str("ErrorCode", constants.MissingParameter).Msg("Missing Container ID")
		response.BadRequestMissingParameters(w, "Missing container ID")
		return
	}
	container, err := h.containerRepo.GetById(containerId)
	if err != nil {
		h.logger.Error().Err(err).Str("ErrorCode", TaskContainerGetNotFound).Msg(err.Error())
		response.NotFound(w, TaskContainerGetNotFound, "Container does not exist")
		return
	}
	response.WriteJsonWithEncode(w, http.StatusOK, container)
}
func (h *Handler) handleGetTaskContainersByGroupId(w http.ResponseWriter, r *http.Request) {
	groupIdVar := chi.URLParam(r, "usergroupID")
	if groupIdVar == "" {
		h.logger.Error().Str("ErrorCode", constants.MissingParameter).Msg("Missing Group ID")
		response.BadRequestMissingParameters(w, "Missing Group ID")
		return
	}
	groupId, err := strconv.Atoi(groupIdVar)
	if err != nil {
		h.logger.Error().Err(err).Str("ErrorCode", constants.InvalidParameter).Msg(err.Error())
		response.ErrorResponse(w, http.StatusBadRequest, *(response.New(constants.InvalidParameter, "Invalid Parameter", "Invalid Group ID")))
		return
	}

	containers, err := h.containerRepo.GetContainersByGroupId(groupId)
	if err != nil {
		h.logger.Error().Err(err).Str("ErrorCode", TaskContainerGetNotFound).Msg(err.Error())
		response.NotFound(w, TaskContainerGetNotFound, "Error occurred during retrieving containers by group id")
		return
	}
	response.WriteJsonWithEncode(w, http.StatusOK, containers)
}
func (h *Handler) handleCreateTaskContainer(w http.ResponseWriter, r *http.Request) {
	var createDto CreateContainerDto
	if err := response.ParseJson(r, &createDto); err != nil {
		h.logger.Error().Err(err).Str("ErrorCode", constants.RequestBodyError).Msg("Error occurred during parsing json of CreateContainerDto")
		response.InvalidJsonBody(w, "Error occurred during parsing json of CreateContainerDto")
		return
	}

	container := model.TaskContainer{
		Id:             uuid.New().String(),
		Name:           createDto.Name,
		Description:    createDto.Description,
		Type:           createDto.Type,
		IsActive:       true,
		Activity_level: 0,
		UsergroupId:    createDto.UserGroupId,
	}
	_ = h.containerRepo.CreateContainer(container)
	response.WriteJsonWithEncode(w, http.StatusCreated, container.Id)
}
func (h *Handler) handleDeleteTaskContainer(w http.ResponseWriter, r *http.Request) {
	containerId := chi.URLParam(r, "containerID")
	err := h.containerRepo.DeleteContainer(containerId)
	if err != nil {
		h.logger.Error().Err(err).Str("ErrorCode", DeleteTaskContainerError).Msg(err.Error())
		response.NotFound(w, DeleteTaskContainerError, "Error occurred during delete container")
		return
	}
	response.WriteJsonWithEncode(w, http.StatusNoContent, "task container is removed.")
}
