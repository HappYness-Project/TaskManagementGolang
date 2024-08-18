package user

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/happYness-Project/taskManagementGolang/internal/auth"
	"github.com/happYness-Project/taskManagementGolang/utils"
)

type Handler struct {
	userRepo UserRepository
}

func NewHandler(repo UserRepository) *Handler {
	return &Handler{userRepo: repo}
}

func (h *Handler) RegisterRoutes(router *chi.Mux) {
	router.Get("/api/users", auth.WithJWTAuth(h.handleGetUsers))
	// router.Get("/api/users", auth.WithJWTAuth(h.handleGetUsers))
	router.Get("/api/users/{userID}", auth.WithJWTAuth(h.handleGetUser))
	router.Get("/api/user-groups/{groupID}/users", auth.WithJWTAuth(h.handleGetUsersByGroupId))
}

func (h *Handler) handleGetUsers(w http.ResponseWriter, r *http.Request) {
	emailVars := r.URL.Query().Get("email")
	userVars := r.URL.Query().Get("username")
	if emailVars != "" {
		user, err := h.userRepo.GetUserByEmail(emailVars)
		if err != nil {
			utils.WriteError(w, http.StatusBadRequest, err)
			return
		}
		userJson, _ := json.Marshal(user)
		utils.WriteJSON(w, http.StatusOK, userJson)
		return
	}
	if userVars != "" {
		user, err := h.userRepo.GetUserByUsername(userVars)
		if err != nil {
			utils.WriteError(w, http.StatusBadRequest, err)
			return
		}
		userJson, _ := json.Marshal(user)
		utils.WriteJSON(w, http.StatusOK, userJson)
		return
	}

	users, err := h.userRepo.GetAllUsers()
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}
	userJson, _ := json.Marshal(users)
	utils.WriteJSON(w, http.StatusOK, userJson)
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
	utils.WriteJsonWithEncode(w, http.StatusOK, user)
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
	utils.WriteJsonWithEncode(w, http.StatusOK, user)
}
