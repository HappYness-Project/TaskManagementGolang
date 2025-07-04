package route

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/jwtauth"
	"github.com/happYness-Project/taskManagementGolang/internal/user/model"
	"github.com/happYness-Project/taskManagementGolang/internal/user/repository"
	userRepo "github.com/happYness-Project/taskManagementGolang/internal/user/repository"
	userGroupRepo "github.com/happYness-Project/taskManagementGolang/internal/usergroup/repository"
	"github.com/happYness-Project/taskManagementGolang/pkg/loggers"
	"github.com/happYness-Project/taskManagementGolang/pkg/utils"
)

type Handler struct {
	logger        *loggers.AppLogger
	userRepo      userRepo.UserRepository
	userGroupRepo userGroupRepo.UserGroupRepository
}

func NewHandler(logger *loggers.AppLogger, repo repository.UserRepository, ugRepo userGroupRepo.UserGroupRepository) *Handler {
	return &Handler{logger: logger, userRepo: repo, userGroupRepo: ugRepo}
}

func (h *Handler) RegisterRoutes(router chi.Router) {
	router.Route("/api/users", func(r chi.Router) {
		r.Get("/", h.handleGetUsers)
		r.Post("/", h.handleCreateUser)
		r.Put("/{userID}", h.handleUpdateUser)
		r.Get("/{userID}", h.handleGetUser)
		r.Patch("/{userID}/default-group", h.handleUpdateGroupId)
	})
	router.Get("/api/user-groups/{groupID}/users", h.handleGetUsersByGroupId)

}

func (h *Handler) handleGetUsers(w http.ResponseWriter, r *http.Request) {

	if r.URL.Query().Get("email") != "" {
		h.responseUser(w, "email", r.URL.Query().Get("email"))
		return
	} else if r.URL.Query().Get("username") != "" {
		h.responseUser(w, "username", r.URL.Query().Get("username"))
		return
	}
	users, err := h.userRepo.GetAllUsers()
	if err != nil {
		h.logger.Error().Err(err).Str("ErrorCode", UserGetServerError).
			Msg("Error occurred during GetAllUsers.")
		utils.ErrorJson(w, err, http.StatusInternalServerError)
		return
	}
	utils.SuccessJson(w, users, "success", http.StatusOK)
}
func (h *Handler) handleGetUser(w http.ResponseWriter, r *http.Request) {
	user, err := h.userRepo.GetUserByUserId(chi.URLParam(r, "userID"))
	if err != nil {
		h.logger.Error().Err(err).Str("ErrorCode", UserGetServerError).Msg("Error occurred during retrieving user.")
		utils.ErrorJson(w, err, http.StatusInternalServerError)
		return
	}
	if user == nil {
		h.logger.Error().Err(err).Str("ErrorCode", UserGetNotFound)
		utils.ErrorJson(w, fmt.Errorf("user does not exist"), http.StatusNotFound)
		return
	}

	userDetailDto := new(UserDetailDto)
	ugs, err := h.userGroupRepo.GetUserGroupsByUserId(user.Id)
	if err != nil {
		h.logger.Error().Err(err).Msg("Bad Request")
		utils.ErrorJson(w, err, http.StatusBadRequest)
		return
	}
	userDetailDto.Id = user.Id
	userDetailDto.UserId = user.UserId
	userDetailDto.UserName = user.UserName
	userDetailDto.FirstName = user.FirstName
	userDetailDto.LastName = user.LastName
	userDetailDto.CreatedAt = user.CreatedAt
	userDetailDto.UpdatedAt = user.UpdatedAt
	userDetailDto.Email = user.Email
	userDetailDto.IsActive = user.IsActive
	userDetailDto.UserGroup = ugs
	userDetailDto.DefaultGroupId = user.DefaultGroupId

	utils.SuccessJson(w, userDetailDto, "success", http.StatusOK)
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
	utils.SuccessJson(w, user, "success", http.StatusOK)
}
func (h *Handler) responseUser(w http.ResponseWriter, findField string, findVar string) {
	var user *model.User
	var err error
	if findField == "email" {
		user, err = h.userRepo.GetUserByEmail(findVar)
	} else if findField == "username" {
		user, err = h.userRepo.GetUserByUsername(findVar)
	}
	if err != nil {
		utils.ErrorJson(w, err, http.StatusBadRequest)
		return
	}
	if user.Id == 0 {
		utils.ErrorJson(w, fmt.Errorf("cannot find user"), http.StatusNotFound)
		return
	}

	userDetailDto := new(UserDetailDto)
	ugs, err := h.userGroupRepo.GetUserGroupsByUserId(user.Id)
	if err != nil {
		utils.ErrorJson(w, err, http.StatusBadRequest)
		return
	}
	userDetailDto.Id = user.Id
	userDetailDto.UserId = user.UserId
	userDetailDto.UserName = user.UserName
	userDetailDto.FirstName = user.FirstName
	userDetailDto.LastName = user.LastName
	userDetailDto.CreatedAt = user.CreatedAt
	userDetailDto.UpdatedAt = user.UpdatedAt
	userDetailDto.Email = user.Email
	userDetailDto.IsActive = user.IsActive
	userDetailDto.UserGroup = ugs
	userDetailDto.DefaultGroupId = user.DefaultGroupId

	utils.SuccessJson(w, userDetailDto, "success", http.StatusOK)
}
func (h *Handler) handleCreateUser(w http.ResponseWriter, r *http.Request) {
	var createDto CreateUserDto
	if err := utils.ParseJson(r, &createDto); err != nil {
		utils.ErrorJson(w, err, http.StatusBadRequest)
		return
	}
	_, claims, _ := jwtauth.FromContext(r.Context())
	userId := fmt.Sprintf("%v", claims["nameid"])

	user := model.NewUser(userId, createDto.UserName, createDto.FirstName, createDto.LastName, createDto.Email)

	err := h.userRepo.CreateUser(*user)
	if err != nil {
		utils.ErrorJson(w, fmt.Errorf("cannot create a user"), http.StatusBadRequest)
		return
	}
	utils.SuccessJson(w, map[string]string{"user_id": userId}, "User is created.", http.StatusCreated)
}
func (h *Handler) handleUpdateUser(w http.ResponseWriter, r *http.Request) {
	var updateDto UpdateUserDto
	if err := utils.ParseJson(r, &updateDto); err != nil {
		utils.ErrorJson(w, err, http.StatusBadRequest)
		return
	}

	user, err := h.userRepo.GetUserByUserId(chi.URLParam(r, "userID"))
	if err != nil || user == nil {
		utils.ErrorJson(w, fmt.Errorf("cannot find user"), http.StatusNotFound)
		return
	}

	user.UpdateUser(updateDto.FirstName, updateDto.LastName, updateDto.Email)

	err = h.userRepo.UpdateUser(*user)
	if err != nil {
		utils.ErrorJson(w, err, http.StatusBadRequest)
		return
	}
	utils.SuccessJson(w, nil, "User is updated.", http.StatusNoContent)
}

func (h *Handler) handleUpdateGroupId(w http.ResponseWriter, r *http.Request) {
	user, err := h.userRepo.GetUserByUserId(chi.URLParam(r, "userID"))
	if err != nil || user == nil {
		utils.ErrorJson(w, fmt.Errorf("cannot find user"), http.StatusNotFound)
		return
	}

	type JsonBody struct {
		DefaultGroupId int `json:"default_group_id"`
	}
	var jsonBody JsonBody
	err = json.NewDecoder(r.Body).Decode(&jsonBody)
	if err != nil {
		utils.ErrorJson(w, err, http.StatusBadRequest)
		return
	}

	err = user.UpdateDefaultGroupId(jsonBody.DefaultGroupId)
	if err != nil {
		utils.ErrorJson(w, err, http.StatusBadRequest)
		return
	}

	err = h.userRepo.UpdateUser(*user)
	if err != nil {
		utils.ErrorJson(w, err, http.StatusBadRequest)
		return
	}

	utils.SuccessJson(w, nil, "Default user group ID is updated.", http.StatusOK)
}
