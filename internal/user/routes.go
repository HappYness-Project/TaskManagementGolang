package user

import (
	"net/http"

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
	router.Get("/users", h.handleGetUsers)
	// router.HandleFunc("/users/{userId}", h.handleGetUser).Methods(http.MethodGet)
}

func (h *Handler) handleGetUsers(w http.ResponseWriter, r *http.Request) {
	users, err := h.userRepo.GetAllUsers()
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}
	utils.WriteJSON(w, http.StatusOK, users)
}

// func (h *Handler) handleGetUser(w http.ResponseWriter, r *http.Request) {
// 	userid := chi.URLParam(r, "userId")
// 	ctx := r.Context()
// 	user, ok := ctx.Value("user").(*User)
// 	if !ok {
// 		http.Error(w, http.StatusText(422), 422)
// 		return
// 	}
// 	w.Write([]byte(fmt.Sprintf("Username:%s", user.UserName)))

// }
