package utils

import (
    "encoding/json"
    "net/http"
)

type ProblemDetails struct {
    Type     string `json:"type,omitempty"`
    Title    string `json:"title"`
    Status   int    `json:"status"`
    Detail   string `json:"detail,omitempty"`
    Instance string `json:"instance,omitempty"`
}

func WriteProblem(w http.ResponseWriter, status int, title, detail, typ, instance string) {
    problem := ProblemDetails{
        Type:     typ,
        Title:    title,
        Status:   status,
        Detail:   detail,
        Instance: instance,
    }
    w.Header().Set("Content-Type", "application/problem+json")
    w.WriteHeader(status)
    _ = json.NewEncoder(w).Encode(problem)
}