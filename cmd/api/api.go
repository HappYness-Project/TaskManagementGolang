package api

import (
	"database/sql"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/jwtauth"

	chatRepo "github.com/happYness-Project/taskManagementGolang/internal/chat/repository"
	taskRepo "github.com/happYness-Project/taskManagementGolang/internal/task/repository"
	containerRepo "github.com/happYness-Project/taskManagementGolang/internal/taskcontainer/repository"
	userRepo "github.com/happYness-Project/taskManagementGolang/internal/user/repository"
	usergroupRepo "github.com/happYness-Project/taskManagementGolang/internal/usergroup/repository"

	chatRoute "github.com/happYness-Project/taskManagementGolang/internal/chat/route"
	taskRoute "github.com/happYness-Project/taskManagementGolang/internal/task/route"
	containerRoute "github.com/happYness-Project/taskManagementGolang/internal/taskcontainer/route"
	userRoute "github.com/happYness-Project/taskManagementGolang/internal/user/route"
	usergroupRoute "github.com/happYness-Project/taskManagementGolang/internal/usergroup/route"
	"github.com/happYness-Project/taskManagementGolang/pkg/configs"
	"github.com/happYness-Project/taskManagementGolang/pkg/loggers"
	"github.com/happYness-Project/taskManagementGolang/pkg/middlewares"
	"github.com/happYness-Project/taskManagementGolang/pkg/response"
)

type ApiServer struct {
	addr      string
	db        *sql.DB
	tokenAuth *jwtauth.JWTAuth
	logger    *loggers.AppLogger
}

func NewApiServer(addr string, db *sql.DB, logger *loggers.AppLogger) *ApiServer {
	tokenAuth := jwtauth.New("HS512", []byte(configs.AccessToken), nil)

	return &ApiServer{
		addr:      addr,
		db:        db,
		tokenAuth: tokenAuth,
		logger:    logger,
	}
}

func (s *ApiServer) Setup() *chi.Mux {
	mux := chi.NewRouter()
	mux.Use(middleware.Logger)
	mux.Use(middleware.Recoverer)
	mux.Use(middlewares.EnableCORS)
	mux.Use(middlewares.RequestIdMiddleware)
	mux.Use(middlewares.Logger(s.logger))
	mux.Use(middleware.Heartbeat("/ping"))
	mux.Get("/", home)
	mux.Get("/health", home)

	userRepo := userRepo.NewUserRepository(s.db)
	usergroupRepo := usergroupRepo.NewUserGroupRepository(s.db)
	taskRepo := taskRepo.NewTaskRepository(s.db)
	containerRepo := containerRepo.NewContainerRepository(s.db)
	chatRepo := chatRepo.NewChatRepository(s.db)

	userHandler := userRoute.NewHandler(s.logger, userRepo, usergroupRepo)
	usergroupHandler := usergroupRoute.NewHandler(s.logger, usergroupRepo, userRepo)
	taskHandler := taskRoute.NewHandler(s.logger, taskRepo, containerRepo, usergroupRepo)
	containerHandler := containerRoute.NewHandler(s.logger, containerRepo, userRepo)
	chatHandler := chatRoute.NewHandler(s.logger, chatRepo, usergroupRepo)

	mux.Group(func(r chi.Router) {
		r.Use(jwtauth.Verifier(s.tokenAuth))
		r.Use(jwtauth.Authenticator)
		userHandler.RegisterRoutes(r)
		usergroupHandler.RegisterRoutes(r)
		taskHandler.RegisterRoutes(r)
		containerHandler.RegisterRoutes(r)
		chatHandler.RegisterRoutes(r)
	})

	return mux
}

func (s *ApiServer) Run(mux *chi.Mux) error {
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
	response.WriteJsonWithEncode(w, http.StatusOK, payload)
}
