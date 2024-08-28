package usergroup

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/happYness-Project/taskManagementGolang/internal/auth"
	"github.com/happYness-Project/taskManagementGolang/utils"
)

type Handler struct {
	groupRepo UserGroupRepository
}

func NewHandler(repo UserGroupRepository) *Handler {
	return &Handler{groupRepo: repo}
}
func (h *Handler) RegisterRoutes(router *chi.Mux) {
	router.Route("/api/user-groups", func(r chi.Router) {
		r.Get("/", auth.WithJWTAuth(h.handleGetUserGroups))
		r.Get("/{groupID}", auth.WithJWTAuth(h.handleGetUserGroupById))
	})
	router.Post("/api/users/{userID}/user-groups", auth.WithJWTAuth(h.handleCraeteUserGroup))

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
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("invalid user ID"))
		return
	}

	group, err := h.groupRepo.GetById(groupId)
	if err != nil {
		utils.WriteError(w, http.StatusNotFound, fmt.Errorf("group does not exist"))
		return
	}
	utils.WriteJsonWithEncode(w, http.StatusOK, group)
}

func (h *Handler) handleCraeteUserGroup(w http.ResponseWriter, r *http.Request) {
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

	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Error reading request body", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	var createDto *CreateUserGroupDto
	err = json.Unmarshal(body, &createDto)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
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
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}

	err = h.groupRepo.InsertUserGroupUserTable(groupId, userId)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}

	utils.WriteJsonWithEncode(w, http.StatusCreated, "User group is created.")

}
