package main

import (
	"fmt"
	"net/http"

	"example.com/taskapp/internal/models"
)

func (app *application) Home(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello World from %s", app.Domain)
}
func (app *application) allTaskContainers(w http.ResponseWriter, r *http.Request) {
	var TaskContainers []models.TaskContainer
	container := models.TaskContainer{ContainerId: "containerIdTestikng", ContainerName: "", ContainerDesc: ""}

	TaskContainers = append(TaskContainers, container)

	fmt.Println("Endpoint hit: All Task Containers.")
}
