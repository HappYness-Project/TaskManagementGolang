package route

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/jwtauth"
	userRepo "github.com/happYness-Project/taskManagementGolang/internal/user/repository"
	userGroupRepo "github.com/happYness-Project/taskManagementGolang/internal/usergroup/repository"

	"github.com/happYness-Project/taskManagementGolang/internal/usergroup/model"
	"github.com/happYness-Project/taskManagementGolang/internal/usergroup/repository"
	"github.com/happYness-Project/taskManagementGolang/pkg/constants"
	"github.com/happYness-Project/taskManagementGolang/pkg/errors"
	"github.com/happYness-Project/taskManagementGolang/pkg/loggers"
	"github.com/happYness-Project/taskManagementGolang/pkg/response"
)

type Handler struct {
	logger    *loggers.AppLogger
	groupRepo userGroupRepo.UserGroupRepository
	userRepo  userRepo.UserRepository
}

func NewHandler(logger *loggers.AppLogger, repo repository.UserGroupRepository, userRepo userRepo.UserRepository) *Handler {
	return &Handler{logger: logger, groupRepo: repo, userRepo: userRepo}
}
func (h *Handler) RegisterRoutes(router chi.Router) {
	router.Route("/api/user-groups", func(r chi.Router) {
		r.Get("/", h.handleGetUserGroups)
		r.Get("/{groupID}", h.handleGetUserGroupById)
		r.Delete("/{groupID}", h.handleDeleteUserGroup)
		r.Post("/{groupID}/users", h.handleAddUserToGroup)
		r.Put("/{groupID}/users/{userID}", h.handleRemoveUserFromGroup)
	})
	router.Post("/api/user-groups", h.handleCreateUserGroup)
	router.Get("/api/users/{userID}/user-groups", h.handleGetUserGroupByUserId)

}

func (h *Handler) handleGetUserGroups(w http.ResponseWriter, r *http.Request) {
	groups, err := h.groupRepo.GetAllUsergroups()
	if err != nil {
		h.logger.Error().Err(err).Str("ErrorCode", constants.ServerError).Msg("Error occurred during responseUser.")
		response.InternalServerError(w)
		return
	}
	response.WriteJsonWithEncode(w, http.StatusOK, groups) // TODO this response format needs to be changed.
}
func (h *Handler) handleGetUserGroupById(w http.ResponseWriter, r *http.Request) {
	groupId, err := strconv.Atoi(chi.URLParam(r, "groupID"))
	if err != nil {
		h.logger.Error().Err(err).Str("ErrorCode", constants.InvalidParameter).Msg("Invalid Parameters for GroupID")
		response.ErrorResponse(w, http.StatusBadRequest, *response.New(constants.InvalidParameter, "Invalid Paramters", "Invalid Group ID"))
		return
	}

	group, err := h.groupRepo.GetById(groupId)
	if group.GroupId == 0 {
		h.logger.Error().Str("ErrorCode", UserGroupGetNotFound).Msg(fmt.Sprintf("group does not exist. group Id: %d", groupId))
		response.NotFound(w, UserGroupGetNotFound, "group does not exist")
		return
	}

	if err != nil {
		h.logger.Error().Err(err).Str("ErrorCode", UserGroupServerError).Msg(err.Error())
		response.ErrorResponse(w, http.StatusBadRequest, *response.New(UserGroupServerError, "Bad Request", "Error occurred during retrieving by group id"))
		return
	}
	response.WriteJsonWithEncode(w, http.StatusOK, group)
}

// TODO List Testing this method
func (h *Handler) handleCreateUserGroup(w http.ResponseWriter, r *http.Request) {
	var createDto CreateUserGroupDto
	if err := response.ParseJson(r, &createDto); err != nil {
		h.logger.Error().Err(err).Str("ErrorCode", constants.RequestBodyError).Msg("Json body for CreateUserGroupRequest is invalid")
		response.InternalServerError(w)
		return
	}
	_, claims, _ := jwtauth.FromContext(r.Context())
	userid := fmt.Sprintf("%v", claims["nameid"])
	user, err := h.userRepo.GetUserByUserId(userid)
	if err != nil || user == nil {
		h.logger.Error().Err(err).Str("ErrorCode", UserGroupGetNotFound).Msg("Not able to find user ID:" + userid)
		response.ErrorResponse(w, http.StatusNotFound, *(response.New(constants.ServerError, errors.InternalServerError)))
		return
	}

	group, err := model.NewUserGroup(createDto.GroupName, createDto.GroupDesc, createDto.GroupType)
	if err != nil {
		h.logger.Error().Err(err).Str("ErrorCode", UserGroupDomainError).Msg(err.Error())
		response.ErrorResponse(w, http.StatusUnprocessableEntity, *response.New(UserGroupDomainError, "Domain Validation Error", err.Error()))
		return
	}

	groupId, err := h.groupRepo.CreateGroupWithUsers(*group, user.Id)
	if err != nil {
		h.logger.Error().Err(err).Str("ErrorCode", UserGroupCreationFailure).Msg(err.Error())
		response.ErrorResponse(w, http.StatusBadRequest, *response.New(UserGroupCreationFailure, "Inserting usergroup failed"))
		return
	}

	response.SuccessJson(w, map[string]int{"group_id": groupId}, "User group is created.", http.StatusCreated)
}

func (h *Handler) handleGetUserGroupByUserId(w http.ResponseWriter, r *http.Request) {
	userId := chi.URLParam(r, "userID")
	if userId == "" {
		h.logger.Error().Msg("missing userID")
		response.BadRequestMissingParameters(w)
		return
	}

	user, err := h.userRepo.GetUserByUserId(userId)
	if err != nil {
		h.logger.Error().Err(err).Str("ErrorCode", UserNotFound).Msg("Error during GetUserByUserId")
		response.NotFound(w, UserNotFound, "Not able to find a user")
		return
	}

	groups, err := h.groupRepo.GetUserGroupsByUserId(user.Id)
	if err != nil {
		h.logger.Error().Err(err).Str("ErrorCode", UserGroupGetNotFound).Msg("Error during GetUserGroupsByUserId")
		response.NotFound(w, UserGroupGetNotFound, "Not able to find usegroups")
		return
	}
	response.WriteJsonWithEncode(w, http.StatusOK, groups)
}

func (h *Handler) handleAddUserToGroup(w http.ResponseWriter, r *http.Request) {
	groupId, err := strconv.Atoi(chi.URLParam(r, "groupID"))
	if err != nil {
		h.logger.Error().Err(err).Str("ErrorCode", constants.InvalidParameter).Msg("invalid Group Id")
		response.BadRequestMissingParameters(w)
		return
	}

	type JsonBody struct {
		UserId string `json:"user_id"`
	}
	var jsonBody JsonBody
	err = json.NewDecoder(r.Body).Decode(&jsonBody)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		h.logger.Error().Err(err).Str("ErrorCode", constants.RequestBodyError)
		response.InvalidJsonBody(w, err.Error())
		return
	}
	user, err := h.userRepo.GetUserByUserId(jsonBody.UserId)
	if err != nil || user == nil {
		h.logger.Error().Err(err).Str("ErrorCode", UserNotFound)
		response.NotFound(w, UserNotFound, "cannot find an user")
		return
	}

	err = h.groupRepo.InsertUserGroupUserTable(groupId, user.Id)
	if err != nil {
		h.logger.Error().Err(err).Str("ErrorCode", UserGroupAddUserError).Msg("Error during InsertUserGroupUserTable")
		response.ErrorResponse(w, http.StatusBadRequest, *response.New(UserGroupAddUserError, "Bad Request", "Inserting usergroup failed"))
		return
	}

	response.WriteJsonWithEncode(w, http.StatusCreated, fmt.Sprintf("User is added to the user group ID: %d", groupId))
}

func (h *Handler) handleDeleteUserGroup(w http.ResponseWriter, r *http.Request) {
	groupId, err := strconv.Atoi(chi.URLParam(r, "groupID"))
	if err != nil {
		h.logger.Error().Err(err).Str("ErrorCode", constants.InvalidParameter).Msg("Invalid groupId")
		response.BadRequestMissingParameters(w, "invalid groupId")
		return
	}
	err = h.groupRepo.DeleteUserGroup(groupId)
	if err != nil {
		h.logger.Error().Err(err).Str("ErrorCode", DeleteUserGroupError).Msg(err.Error())
		// TODO : Possible error can occur if not found.
		response.ErrorResponse(w, http.StatusBadRequest, *(response.New(DeleteUserGroupError, "Bad Request", "Error occurred during deleting user group.")))
		return
	}
	response.SuccessJson(w, nil, fmt.Sprintf("User is removed from user group ID: %d", groupId), 204)
}

func (h *Handler) handleRemoveUserFromGroup(w http.ResponseWriter, r *http.Request) {
	vars := chi.URLParam(r, "groupID")
	if vars == "" {
		h.logger.Error().Str("ErrorCode", constants.MissingParameter).Msg("Missing GroupID")
		response.BadRequestMissingParameters(w, "Missing Group ID")
		return
	}
	groupId, err := strconv.Atoi(vars)
	if err != nil {
		h.logger.Error().Err(err).Str("ErrorCode", constants.InvalidParameter).Msg(err.Error())
		response.ErrorResponse(w, http.StatusBadRequest, *(response.New(constants.InvalidParameter, "Invalid Parameter", "Invalid Group ID")))
		return
	}
	userId := chi.URLParam(r, "userID")
	user, err := h.userRepo.GetUserByUserId(userId)
	if err != nil || user == nil {
		h.logger.Error().Err(err).Str("ErrorCode", UserNotFound).Msg(err.Error())
		response.NotFound(w, UserNotFound, "Provided user id cannot be found")
		return
	}

	err = h.groupRepo.RemoveUserFromUserGroup(groupId, user.Id)
	if err != nil {
		h.logger.Error().Err(err).Str("ErrorCode", RemoveUserFromUserGroupError).Msg(err.Error())
		response.ErrorResponse(w, http.StatusBadRequest, *(response.New(RemoveUserFromUserGroupError, "Bad Request", "Failed to remove a user from usergroup.")))
		return
	}
	response.SuccessJson(w, nil, fmt.Sprintf("User is removed from user group ID: %d", groupId), 204)
}
