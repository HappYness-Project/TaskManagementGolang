package usergroup

import (
	"fmt"
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
