package user

import (
	"fmt"
	"net/http"
	"strconv"
	"time"

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
	router.Get("/api/users", h.handleGetUsers)
	router.Post("/api/users", h.handleCreateUser)
	router.Get("/api/users/{userID}", h.handleGetUser)
	router.Get("/api/user-groups/{groupID}/users", h.handleGetUsersByGroupId)
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
	utils.WriteJsonWithEncode(w, http.StatusOK, users)
}
func (h *Handler) handleGetUser(w http.ResponseWriter, r *http.Request) {
	userID, err := strconv.Atoi(chi.URLParam(r, "userID"))
	if err != nil {
		utils.ErrorJSON(w, fmt.Errorf("invalid user ID"), http.StatusBadRequest)
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
		utils.ErrorJSON(w, fmt.Errorf("missing Group ID"), http.StatusBadRequest)
		return
	}
	groupId, err := strconv.Atoi(vars)
	if err != nil {
		utils.ErrorJSON(w, fmt.Errorf("invalid user ID"), http.StatusBadRequest)
		return
	}
	user, err := h.userRepo.GetUsersByGroupId(groupId)
	if err != nil {
		utils.ErrorJSON(w, fmt.Errorf("user does not exist"), http.StatusNotFound)
		return
	}
	utils.WriteJsonWithEncode(w, http.StatusOK, user)
}
func (h *Handler) responseUserUsingEmail(w http.ResponseWriter, findField string, findVar string) {
	var user *User
	var err error
	if findField == "email" {
		user, err = h.userRepo.GetUserByEmail(findVar)
	} else if findField == "username" {
		user, err = h.userRepo.GetUserByUsername(findVar)
	}
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}
	if user.Id == 0 {
		utils.WriteError(w, http.StatusNotFound, fmt.Errorf("cannot find user"))
		return
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
	userDetailDto.CreatedAt = user.CreatedAt
	userDetailDto.UpdatedAt = user.UpdatedAt
	userDetailDto.Email = user.Email
	userDetailDto.IsActive = user.IsActive
	userDetailDto.UserGroup = ugs
	usersetting, _ := h.userRepo.GetGroupSettingByUserId(user.Id)
	userDetailDto.UserSetting = usersetting

	utils.WriteJsonWithEncode(w, http.StatusOK, userDetailDto)
}
func (h *Handler) handleCreateUser(w http.ResponseWriter, r *http.Request) {
	var createDto CreateUserDto
	if err := utils.ParseJSON(r, &createDto); err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}

	user := User{
		Id:        auth.GetUserIDFromContext(r.Context()),
		UserName:  createDto.UserName,
		FirstName: createDto.FirstName,
		LastName:  createDto.LastName,
		Email:     createDto.Email,
		IsActive:  true,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	err := h.userRepo.Create(user)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}
	utils.WriteJsonWithEncode(w, http.StatusCreated, "User is created.")
}
