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
	"github.com/happYness-Project/taskManagementGolang/pkg/loggers"
	"github.com/happYness-Project/taskManagementGolang/pkg/utils"
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
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}
	utils.WriteJsonWithEncode(w, http.StatusOK, groups)
}
func (h *Handler) handleGetUserGroupById(w http.ResponseWriter, r *http.Request) {
	vars := chi.URLParam(r, "groupID")
	if vars == "" {
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("missing Group ID"))
		return
	}
	groupId, err := strconv.Atoi(vars)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("invalid group ID"))
		return
	}

	group, err := h.groupRepo.GetById(groupId)
	if err != nil {
		utils.WriteError(w, http.StatusNotFound, fmt.Errorf("group does not exist"))
		return
	}
	utils.WriteJsonWithEncode(w, http.StatusOK, group)
}

// TODO List Testing this method
func (h *Handler) handleCreateUserGroup(w http.ResponseWriter, r *http.Request) {
	var createDto CreateUserGroupDto
	if err := utils.ParseJson(r, &createDto); err != nil {
		utils.ErrorJson(w, fmt.Errorf("error reading request body"), http.StatusBadRequest)
		return
	}
	_, claims, _ := jwtauth.FromContext(r.Context())
	userid := fmt.Sprintf("%v", claims["nameid"])
	user, err := h.userRepo.GetUserByUserId(userid)
	if err != nil || user == nil {
		utils.ErrorJson(w, fmt.Errorf("cannot find user"), http.StatusNotFound)
		return
	}

	group, err := model.NewUserGroup(createDto.GroupName, createDto.GroupDesc, createDto.GroupType)
	if err != nil {
		utils.ErrorJson(w, err, http.StatusBadRequest)
		return
	}

	groupId, err := h.groupRepo.CreateGroupWithUsers(*group, user.Id)
	if err != nil {

		utils.ErrorJson(w, err, http.StatusBadRequest)
		return
	}

	utils.SuccessJson(w, map[string]int{"group_id": groupId}, "User group is created.", http.StatusCreated)
}

func (h *Handler) handleGetUserGroupByUserId(w http.ResponseWriter, r *http.Request) {
	userId := chi.URLParam(r, "userID")
	if userId == "" {
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("missing user ID"))
		return
	}

	user, err := h.userRepo.GetUserByUserId(userId)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("user does not exist."))
	}

	groups, err := h.groupRepo.GetUserGroupsByUserId(user.Id)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("error occurred during GetUserGroupsByUserId"))
		return
	}
	utils.WriteJsonWithEncode(w, http.StatusOK, groups)
}

func (h *Handler) handleAddUserToGroup(w http.ResponseWriter, r *http.Request) {
	vars := chi.URLParam(r, "groupID")
	if vars == "" {
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("missing Group ID"))
		return
	}
	groupId, err := strconv.Atoi(vars)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("invalid group ID"))
		return
	}

	type JsonBody struct {
		UserId string `json:"user_id"`
	}
	var jsonBody JsonBody
	err = json.NewDecoder(r.Body).Decode(&jsonBody)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	user, err := h.userRepo.GetUserByUserId(jsonBody.UserId)
	if err != nil || user == nil {
		utils.ErrorJson(w, fmt.Errorf("cannot find user"), http.StatusNotFound)
		return
	}

	err = h.groupRepo.InsertUserGroupUserTable(groupId, user.Id)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}

	utils.WriteJsonWithEncode(w, http.StatusCreated, fmt.Sprintf("User is added to the user group ID: %d", groupId))
}

func (h *Handler) handleDeleteUserGroup(w http.ResponseWriter, r *http.Request) {
	groupId, err := strconv.Atoi(chi.URLParam(r, "groupID"))
	if err != nil {
		utils.ErrorJson(w, fmt.Errorf("invalid Group ID"), http.StatusBadRequest)
		return
	}
	err = h.groupRepo.DeleteUserGroup(groupId)
	if err != nil {
		utils.ErrorJson(w, err, http.StatusBadRequest)
		return
	}
	utils.SuccessJson(w, nil, fmt.Sprintf("User is removed from user group ID: %d", groupId), 204)
}

func (h *Handler) handleRemoveUserFromGroup(w http.ResponseWriter, r *http.Request) {
	vars := chi.URLParam(r, "groupID")
	if vars == "" {
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("missing Group ID"))
		return
	}
	groupId, err := strconv.Atoi(vars)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("invalid Group ID"))
		return
	}
	userId := chi.URLParam(r, "userID")
	user, err := h.userRepo.GetUserByUserId(userId)
	if err != nil || user == nil {
		utils.ErrorJson(w, fmt.Errorf("cannot find user"), http.StatusNotFound)
		return
	}

	err = h.groupRepo.RemoveUserFromUserGroup(groupId, user.Id)
	if err != nil {
		utils.ErrorJson(w, fmt.Errorf("failed to remove user to the user group"), 400)
		return
	}
	utils.SuccessJson(w, nil, fmt.Sprintf("User is removed from user group ID: %d", groupId), 204)
}
