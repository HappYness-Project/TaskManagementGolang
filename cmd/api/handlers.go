package main

// func (app *application) Home(w http.ResponseWriter, r *http.Request) {

// 	var payload = struct {
// 		Status  string `json:"status"`
// 		Message string `json:"message"`
// 		Version string `json:"version"`
// 	}{
// 		Status:  "active",
// 		Message: "Golang Task manager app",
// 		Version: "1.0.0",
// 	}
// 	_ = app.WriteJson(w, http.StatusOK, payload)
// }
// func (app *application) allTaskContainers(w http.ResponseWriter, r *http.Request) {
// 	containers, err := app.containerRepo.AllTaskContainers()
// 	if err != nil {
// 		app.errorJson(w, err)
// 		return
// 	}
// 	_ = app.WriteJson(w, http.StatusOK, containers)
// }
