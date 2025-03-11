package route

import (
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/jwtauth"
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
	})
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
		h.logger.Error().Err(err).Str("ErrorCode", UserGetNotFound).
			Msg("Error occurred during GetUserById.")
		utils.ErrorJson(w, fmt.Errorf("user does not exist"), http.StatusNotFound)
		return
	}
	utils.SuccessJson(w, user, "success", http.StatusOK)
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

func (h *Handler) responseUserUsingEmail(w http.ResponseWriter, findField string, findVar string) {
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

	user := model.User{
		UserId:         userId,
		UserName:       createDto.UserName,
		FirstName:      createDto.FirstName,
		LastName:       createDto.LastName,
		Email:          createDto.Email,
		IsActive:       true,
		CreatedAt:      time.Now(),
		UpdatedAt:      time.Now(),
		DefaultGroupId: 0,
	}

	err := h.userRepo.CreateUser(user)
	if err != nil {
		utils.ErrorJson(w, fmt.Errorf("cannot create a user"), http.StatusBadRequest)
		return
	}
	utils.SuccessJson(w, nil, "User is created.", http.StatusCreated)
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
	updatedUser := model.User{
		Id:        user.Id,
		UserId:    user.UserId,
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
		utils.ErrorJson(w, err, http.StatusBadRequest)
		return
	}
	// TODO no success message is displayed.
	utils.SuccessJson(w, nil, "User is updated.", http.StatusNoContent)

}
