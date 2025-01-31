package usergroup

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/jwtauth"
	"github.com/happYness-Project/taskManagementGolang/pkg/utils"
)

type Handler struct {
	groupRepo UserGroupRepository
}

func NewHandler(repo UserGroupRepository) *Handler {
	return &Handler{groupRepo: repo}
}
func (h *Handler) RegisterRoutes(router *chi.Mux) {
	router.Route("/api/user-groups", func(r chi.Router) {
		r.Get("/", h.handleGetUserGroups)
		r.Get("/{groupID}", h.handleGetUserGroupById)
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

func (h *Handler) handleCreateUserGroup(w http.ResponseWriter, r *http.Request) {
	_, claims, _ := jwtauth.FromContext(r.Context())
	userID, err := strconv.Atoi(fmt.Sprintf("%v", claims["nameid"]))
	if err != nil {
		utils.ErrorJson(w, fmt.Errorf("invalid user ID"), http.StatusBadRequest)
		return
	}
	var createDto *CreateUserGroupDto
	if err := utils.ParseJson(r, &createDto); err != nil {
		utils.ErrorJson(w, fmt.Errorf("error reading request body"), http.StatusBadRequest)
		return
	}

	err = ValidateNewUserGroup(*createDto)
	if err != nil {
		utils.ErrorJson(w, err, http.StatusBadRequest)
		return
	}

	group := UserGroup{
		GroupName: createDto.GroupName,
		GroupDesc: createDto.GroupDesc,
		Type:      createDto.GroupType,
		IsActive:  true,
		Thumbnail: "",
	}
	groupId, err := h.groupRepo.CreateGroup(group)
	if err != nil {
		utils.ErrorJson(w, err, http.StatusBadRequest)
		return
	}
	err = h.groupRepo.InsertUserGroupUserTable(groupId, userID)
	if err != nil {
		utils.ErrorJson(w, err, http.StatusBadRequest)
		return
	}

	utils.SuccessJson(w, nil, "User group is created.", http.StatusCreated)
}

func ValidateNewUserGroup(req CreateUserGroupDto) error {
	if req.GroupName == "" {
		return fmt.Errorf("GroupName field cannot be empty")
	}
	if req.GroupType == "" {
		return fmt.Errorf("GroupType field cannot be empty")
	}
	return nil
}

func (h *Handler) handleGetUserGroupByUserId(w http.ResponseWriter, r *http.Request) {
	vars := chi.URLParam(r, "userID")
	if vars == "" {
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("missing user ID"))
		return
	}
	userId, err := strconv.Atoi(vars)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("invalid user ID"))
		return
	}

	groups, err := h.groupRepo.GetUserGroupsByUserId(userId)
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
		UserId int `json:"user_id"`
	}
	var jsonBody JsonBody
	err = json.NewDecoder(r.Body).Decode(&jsonBody)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = h.groupRepo.InsertUserGroupUserTable(groupId, jsonBody.UserId)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}

	utils.WriteJsonWithEncode(w, http.StatusCreated, fmt.Sprintf("User is added to the user group ID: %d", groupId))
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
	vars = chi.URLParam(r, "userID")
	if vars == "" {
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("missing User ID"))
		return
	}
	userId, err := strconv.Atoi(vars)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("invalid User ID"))
		return
	}

	err = h.groupRepo.RemoveUserFromUserGroup(groupId, userId)
	if err != nil {
		utils.ErrorJson(w, fmt.Errorf("failed to remove user to the user group"), 400)
		return
	}
	utils.SuccessJson(w, nil, fmt.Sprintf("User is removed from user group ID: %d", 204))
}
