package route

import (
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/happYness-Project/taskManagementGolang/internal/auth"
	"github.com/happYness-Project/taskManagementGolang/internal/user/model"
	"github.com/happYness-Project/taskManagementGolang/internal/user/repository"
	"github.com/happYness-Project/taskManagementGolang/internal/usergroup"
	"github.com/happYness-Project/taskManagementGolang/pkg/loggers"
	"github.com/happYness-Project/taskManagementGolang/pkg/utils"
)

type Handler struct {
	logger        *loggers.AppLogger
	userRepo      repository.UserRepository
	userGroupRepo usergroup.UserGroupRepository
}

func NewHandler(logger *loggers.AppLogger, repo repository.UserRepository, ugRepo usergroup.UserGroupRepository) *Handler {
	return &Handler{userRepo: repo, userGroupRepo: ugRepo}
}

func (h *Handler) RegisterRoutes(router *chi.Mux) {

	router.Route("/api/users", func(r chi.Router) {
		r.Get("/", h.handleGetUsers)
		r.Post("/", h.handleCreateUser)
		r.Put("/{userID}", h.handleUpdateUser)
		r.Get("/{userID}", h.handleGetUser)
		r.Get("/{groupID}/users", h.handleGetUsersByGroupId)
	})
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
		h.logger.Error().Err(err).Str("ErrorCode", UserGetServerError).
			Msg("Error occurred during GetAllUsers.")
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}
	utils.WriteJsonWithEncode(w, http.StatusOK, users)
}
func (h *Handler) handleGetUser(w http.ResponseWriter, r *http.Request) {
	userID, err := strconv.Atoi(chi.URLParam(r, "userID"))
	if err != nil {
		utils.ErrorJson(w, fmt.Errorf("invalid user ID"), http.StatusBadRequest)
		return
	}
	user, err := h.userRepo.GetUserById(userID)
	if err != nil {
		h.logger.Error().Err(err).Str("ErrorCode", UserGetNotFound).
			Msg("Error occurred during GetUserById.")
		utils.WriteError(w, http.StatusNotFound, fmt.Errorf("user does not exist"))
		return
	}
	utils.WriteJsonWithEncode(w, http.StatusOK, user)
}
func (h *Handler) handleGetUsersByGroupId(w http.ResponseWriter, r *http.Request) {
	vars := chi.URLParam(r, "groupID")
	if vars == "" {
		utils.ErrorJson(w, fmt.Errorf("missing Group ID"), http.StatusBadRequest)
		return
	}
	groupId, err := strconv.Atoi(vars)
	if err != nil {
		utils.ErrorJson(w, fmt.Errorf("invalid user ID"), http.StatusBadRequest)
		return
	}
	user, err := h.userRepo.GetUsersByGroupId(groupId)
	if err != nil {
		utils.ErrorJson(w, fmt.Errorf("user does not exist"), http.StatusNotFound)
		return
	}
	utils.WriteJsonWithEncode(w, http.StatusOK, user)
}
func (h *Handler) responseUserUsingEmail(w http.ResponseWriter, findField string, findVar string) {
	var user *model.User
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
	if err := utils.ParseJson(r, &createDto); err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}

	user := model.User{
		Id:        auth.GetUserIDFromContext(r.Context()),
		UserName:  createDto.UserName,
		FirstName: createDto.FirstName,
		LastName:  createDto.LastName,
		Email:     createDto.Email,
		IsActive:  true,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	err := h.userRepo.CreateUser(user)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}
	utils.WriteJsonWithEncode(w, http.StatusCreated, "User is created.")
}
func (h *Handler) handleUpdateUser(w http.ResponseWriter, r *http.Request) {
	var updateDto UpdateUserDto
	if err := utils.ParseJson(r, &updateDto); err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}

	userID, err := strconv.Atoi(chi.URLParam(r, "userID"))
	if err != nil {
		utils.ErrorJson(w, fmt.Errorf("invalid user ID"), http.StatusBadRequest)
		return
	}

	user, err := h.userRepo.GetUserById(userID)
	if err != nil || user == nil {
		utils.WriteError(w, http.StatusNotFound, fmt.Errorf("cannot find user"))
		return
	}

	updatedUser := model.User{
		Id:        auth.GetUserIDFromContext(r.Context()),
		UserName:  user.UserName,
		FirstName: updateDto.FirstName,
		LastName:  updateDto.LastName,
		Email:     updateDto.Email,
		IsActive:  true,
		CreatedAt: user.CreatedAt,
		UpdatedAt: time.Now(),
	}
	err = h.userRepo.UpdateUser(updatedUser)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}

}
