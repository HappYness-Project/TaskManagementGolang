package user

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/happYness-Project/taskManagementGolang/internal/auth"
	"github.com/happYness-Project/taskManagementGolang/internal/usergroup"
	"github.com/happYness-Project/taskManagementGolang/utils"
)

type Handler struct {
	userRepo      UserRepository
	userGroupRepo usergroup.UserGroupRepository
}

func NewHandler(repo UserRepository, ugRepo usergroup.UserGroupRepository) *Handler {
	return &Handler{userRepo: repo, userGroupRepo: ugRepo}
}

func (h *Handler) RegisterRoutes(router *chi.Mux) {
	router.Get("/api/users", auth.WithJWTAuth(h.handleGetUsers))
	// router.Get("/api/users", auth.WithJWTAuth(h.handleGetUsers))
	router.Get("/api/users/{userID}", auth.WithJWTAuth(h.handleGetUser))
	router.Get("/api/user-groups/{groupID}/users", auth.WithJWTAuth(h.handleGetUsersByGroupId))
}

func (h *Handler) handleGetUsers(w http.ResponseWriter, r *http.Request) {
	if r.URL.Query().Get("email") != "" {
		h.responseUserUsingEmail(w, "email", r.URL.Query().Get("email"))
		return
	} else if r.URL.Query().Get("username") != "" {
		h.responseUserUsingEmail(w, "username", r.URL.Query().Get("username"))
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
	userJson, _ := json.Marshal(user)
	utils.WriteJSON(w, http.StatusOK, userJson)
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
func (h *Handler) responseUserUsingEmail(w http.ResponseWriter, findField string, findVar string) {
	var user *User
	var err error
	if findField == "email" {
		user, err = h.userRepo.GetUserByEmail(findVar)
		if err != nil {
			utils.WriteError(w, http.StatusBadRequest, err)
			return
		}
	} else if findField == "username" {
		user, err = h.userRepo.GetUserByUsername(findVar)
		if err != nil {
			utils.WriteError(w, http.StatusBadRequest, err)
			return
		}
	}
	var defaultGroupId int
	if user.UserSettingId != 0 {
		defaultGroupId, err = h.userRepo.GetDefaultGroupId(user.UserSettingId)
		if err != nil {
			utils.WriteError(w, http.StatusBadRequest, err)
			return
		}
	}
	userDetailDto := new(UserDetailDto)
	ugs, err := h.userGroupRepo.GetUserGroupsByUserId(user.Id)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}
	userDetailDto.Id = user.Id
	userDetailDto.UserName = user.UserName
	userDetailDto.FirstName = user.FirstName
	userDetailDto.LastName = user.LastName
	userDetailDto.Email = user.Email
	userDetailDto.IsActive = user.IsActive
	userDetailDto.DefaultGroupId = defaultGroupId
	userDetailDto.UserGroup = ugs

	userJson, _ := json.Marshal(userDetailDto)
	utils.WriteJSON(w, http.StatusOK, userJson)
}
