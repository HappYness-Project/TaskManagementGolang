package middlewares

import (
	"context"
	"net/http"

	"github.com/google/uuid"
	"github.com/happYness-Project/taskManagementGolang/pkg/utils"
)

func RequestIdMiddleware(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		rid := r.Header.Get("X-Request-ID")
		if rid == "" {
			rid = uuid.New().String()
		}
		ctx := context.WithValue(r.Context(), utils.ContextKey(utils.RequestIdentifier), rid)
		w.Header().Add("X-Request-ID", rid)

		h.ServeHTTP(w, r.WithContext(ctx))
	})
}
