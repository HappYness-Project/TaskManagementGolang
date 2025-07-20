package response

import (
	"encoding/json"
	"net/http"

	"github.com/happYness-Project/taskManagementGolang/pkg/constants"
	"github.com/happYness-Project/taskManagementGolang/pkg/errors"
)

func ErrorResponse(w http.ResponseWriter, status int, problem ProblemDetails) {
	w.Header().Set("Content-Type", "application/problem+json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(problem)
}

func InternalServerError(w http.ResponseWriter, detail ...string) {
	ErrorResponse(w, http.StatusInternalServerError, *(New(constants.ServerError, errors.InternalServerError, detail...)))
}
func InvalidJsonBody(w http.ResponseWriter, detail ...string) {
	ErrorResponse(w, http.StatusBadRequest, *(New(constants.RequestBodyError, errors.InvalidJsonBody, detail...)))

}

func BadRequestMissingParameters(w http.ResponseWriter, detail ...string) {
	ErrorResponse(w, http.StatusBadRequest, *(New(constants.MissingParameter, errors.Badrequest, detail...)))
}
func BadRequestDomainError(w http.ResponseWriter, err_code string, details ...string) {
	p := New(err_code, "Domain Error", details...)
	ErrorResponse(w, http.StatusBadRequest, *p)
}
func NotFound(w http.ResponseWriter, err_code string, details ...string) {
	p := New(err_code, "Not found", details...)
	ErrorResponse(w, http.StatusNotFound, *p)
}
