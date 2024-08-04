package main

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func (app *application) routes() http.Handler {
	// containerRepo := repository.NewContainerRepo(app.database)
	// containerRepo.
	mux := chi.NewRouter()
	mux.Use(middleware.Recoverer)
	mux.Get("/", app.Home)
	mux.Get("/task-containers", app.allTaskContainers)
	return mux
}
