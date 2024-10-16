package api

import (
	"database/sql"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/happYness-Project/taskManagementGolang/internal/task"
	"github.com/happYness-Project/taskManagementGolang/internal/taskcontainer"
	"github.com/happYness-Project/taskManagementGolang/internal/user"
	"github.com/happYness-Project/taskManagementGolang/internal/usergroup"
	"github.com/happYness-Project/taskManagementGolang/utils"
)

type ApiServer struct {
	addr string
	db   *sql.DB
}

func NewApiServer(addr string, db *sql.DB) *ApiServer {
	return &ApiServer{
		addr: addr,
		db:   db,
	}
}

func (s *ApiServer) Run() error {
	mux := chi.NewRouter()
	mux.Use(middleware.Recoverer)
	// mux.Use(enableCORS)
	// logging by doing mux.Use(middleware.Logger)

	mux.Get("/", home)

	userRepo := user.NewUserRepository(s.db)
	usergroupRepo := usergroup.NewUserGroupRepository(s.db)
	taskRepo := task.NewTaskRepository(s.db)
	containerRepo := taskcontainer.NewContainerRepository(s.db)

	userHandler := user.NewHandler(userRepo, usergroupRepo)
	usergroupHandler := usergroup.NewHandler(usergroupRepo)
	taskHandler := task.NewHandler(taskRepo, containerRepo, usergroupRepo)
	containerHandler := taskcontainer.NewHandler(containerRepo, userRepo)

	userHandler.RegisterRoutes(mux)
	usergroupHandler.RegisterRoutes(mux)
	taskHandler.RegisterRoutes(mux)
	containerHandler.RegisterRoutes(mux)

	log.Println("Listening on ", s.addr)
	return http.ListenAndServe(s.addr, mux)
}
func home(w http.ResponseWriter, r *http.Request) {

	var payload = struct {
		Status  string `json:"status"`
		Message string `json:"message"`
		Version string `json:"version"`
	}{
		Status:  "active",
		Message: "Golang Task manager app",
		Version: "1.0.0",
	}
	utils.WriteJsonWithEncode(w, http.StatusOK, payload)
}
