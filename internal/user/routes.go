package user

import (
	"fmt"
	"net/http"
	"strconv"

	"example.com/taskapp/utils"
	"github.com/go-chi/chi/v5"
)

type Handler struct {
	userRepo UserRepository
}

func NewHandler(repo UserRepository) *Handler {
	return &Handler{userRepo: repo}
}

func (h *Handler) RegisterRoutes(router *chi.Mux) {
	router.Get("/api/users", h.handleGetUsers)
	router.Get("/api/users/{userID}", h.handleGetUser)
	router.Get("/api/user-groups/{groupID}/users", h.handleGetUsersByGroupId)
}

func (h *Handler) handleGetUsers(w http.ResponseWriter, r *http.Request) {
	users, err := h.userRepo.GetAllUsers()
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}
	utils.WriteJSON(w, http.StatusOK, users)
}

func (h *Handler) handleGetUser(w http.ResponseWriter, r *http.Request) {
	vars := chi.URLParam(r, "userID")
	if vars == "" {
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("missing User ID"))
		return
	}
	userID, err := strconv.Atoi(vars)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("invalid user ID"))
		return
	}
	user, err := h.userRepo.GetUserById(userID)
	if err != nil {
		utils.WriteError(w, http.StatusNotFound, fmt.Errorf("user does not exist"))
		return
	}
	utils.WriteJSON(w, http.StatusOK, user)
}

func (h *Handler) handleGetUsersByGroupId(w http.ResponseWriter, r *http.Request) {
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
	user, err := h.userRepo.GetUsersByGroupId(groupId)
	if err != nil {
		utils.WriteError(w, http.StatusNotFound, fmt.Errorf("user does not exist"))
		return
	}
	utils.WriteJSON(w, http.StatusOK, user)
}
