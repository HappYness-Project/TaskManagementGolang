package route

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	chatRepo "github.com/happYness-Project/taskManagementGolang/internal/chat/repository"
	usergroupRepo "github.com/happYness-Project/taskManagementGolang/internal/usergroup/repository"
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
		// r.Get("/{chat}", h.handleGetChatById)
		// r.Delete("/{chatID}", h.handleDeleteChat)
	})
	// router.Get("/api/task-containers/{containerID}/tasks", h.handleGetTasksByContainerId)
	// router.Post("/api/user-groups/{usergroupID}/chats", h.handleCreateChat)
	// router.Get("/api/user-groups/{usergroupID}/tasks", h.handleGetChatByUsergroupId)
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
