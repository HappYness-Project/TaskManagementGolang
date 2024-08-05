package main

import (
	"database/sql"
	"log"
	"net/http"

	"example.com/taskapp/internal/task"
	"example.com/taskapp/internal/taskcontainer"
	"example.com/taskapp/internal/user"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
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

	userRepo := user.NewUserRepository(s.db)
	userHandler := user.NewHandler(userRepo)
	userHandler.RegisterRoutes(mux)

	containerRepo := taskcontainer.NewContainerRepository(s.db)
	containerHandler := taskcontainer.NewHandler(containerRepo)
	containerHandler.RegisterRoutes(mux)

	taskRepo := task.NewTaskRepository(s.db)
	taskHandler := task.NewHandler(taskRepo)
	taskHandler.RegisterRoutes(mux)

	log.Println("Listening on ", s.addr)
	return http.ListenAndServe(s.addr, mux)
}
