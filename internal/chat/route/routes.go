package route

import (
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	chatRepo "github.com/happYness-Project/taskManagementGolang/internal/chat/repository"
	usergroupRepo "github.com/happYness-Project/taskManagementGolang/internal/usergroup/repository"
	"github.com/happYness-Project/taskManagementGolang/pkg/constants"
	"github.com/happYness-Project/taskManagementGolang/pkg/loggers"
	"github.com/happYness-Project/taskManagementGolang/pkg/response"
)

type Handler struct {
	logger    *loggers.AppLogger
	chatRepo  chatRepo.ChatRepository
	groupRepo usergroupRepo.UserGroupRepository
}

func NewHandler(logger *loggers.AppLogger, repo chatRepo.ChatRepository, ugRepo usergroupRepo.UserGroupRepository) *Handler {
	return &Handler{logger: logger, chatRepo: repo, groupRepo: ugRepo}
}
func (h *Handler) RegisterRoutes(router chi.Router) {
	router.Route("/api/chats", func(r chi.Router) {
		r.Get("/", h.handleGetAllChats)
		r.Get("/{id}", h.handleGetChatById)
		// r.Delete("/{chatID}", h.handleDeleteChat)
	})
	router.Get("/api/user-groups/{usergroupID}/chats", h.handleGetChatByUsergroupId)
	// router.Post("/api/user-groups/{usergroupID}/chats", h.handleCreateChat)
}

func (h *Handler) handleGetAllChats(w http.ResponseWriter, r *http.Request) {
	chats, err := h.chatRepo.GetAllChats()
	if err != nil {
		h.logger.Error().Err(err).Str("ErrorCode", ChatGetServerError).Msg(err.Error())
		response.InternalServerError(w, "Error occurred during getting all chats.")
		return
	}
	response.SuccessJson(w, chats, "successfully get chats", http.StatusOK)
}

func (h *Handler) handleGetChatByUsergroupId(w http.ResponseWriter, r *http.Request) {
	groupId, err := strconv.Atoi(chi.URLParam(r, "usergroupID"))
	if err != nil {
		h.logger.Error().Err(err).Msg("invalid Group ID")
		response.ErrorResponse(w, http.StatusBadRequest, *(response.New(constants.InvalidParameter, "Invalid Group ID")))
		return
	}

	chats, err := h.chatRepo.GetChatByUserGroupId(groupId)
	if err != nil {
		h.logger.Error().Err(err).Str("ErrorCode", ChatGetServerError).Msg(err.Error())
		response.InternalServerError(w, "Error occurred during getting all chats.")
		return
	}
	response.SuccessJson(w, chats, "successfully get chats", http.StatusOK)
}

func (h *Handler) handleGetChatById(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if id == "" {
		h.logger.Error().Msg("missing chat ID in request")
		response.BadRequestMissingParameters(w, "Missing Chat ID")
		return
	}

	chat, err := h.chatRepo.GetChatById(id)
	if err != nil {
		h.logger.Error().Err(err).Str("ErrorCode", ChatGetServerError).Msg("failed to get chat by ID")
		response.InternalServerError(w, "Error occurred during getting chat by ID.")
		return
	}
	if chat == nil {
		h.logger.Error().Str("ErrorCode", ChatGetNotFound).Msg("Not able find a chat")
		response.NotFound(w, ChatGetNotFound)
		return
	}

	response.SuccessJson(w, chat, "successfully get chat", http.StatusOK)
}
