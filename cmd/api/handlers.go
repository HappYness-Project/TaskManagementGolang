package main

import (
	"net/http"
)

func (app *application) Home(w http.ResponseWriter, r *http.Request) {

	var payload = struct {
		Status  string `json:"status"`
		Message string `json:"message"`
		Version string `json:"version"`
	}{
		Status:  "active",
		Message: "Golang Task manager app",
		Version: "1.0.0",
	}
	_ = app.writeJson(w, http.StatusOK, payload)
}
func (app *application) allTaskContainers(w http.ResponseWriter, r *http.Request) {
	containers, err := app.DB.AllTaskContainers()
	if err != nil {
		app.errorJson(w, err)
		return
	}
	_ = app.writeJson(w, http.StatusOK, containers)
}

func (app *application) allTasksByContainerId(w http.ResponseWriter, r *http.Request) {
	// tasks, err := app.DB.
}
