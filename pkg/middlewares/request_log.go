package middlewares

import (
	"net/http"
	"time"

	"github.com/happYness-Project/taskManagementGolang/pkg/loggers"
)

type loggingResponseWriter struct {
	http.ResponseWriter
	statusCode int
}

func NewLoggingResponseWriter(w http.ResponseWriter) *loggingResponseWriter {
	return &loggingResponseWriter{w, http.StatusOK}
}

func (lrw *loggingResponseWriter) WriteHeader(code int) {
	lrw.statusCode = code
	lrw.ResponseWriter.WriteHeader(code)
}

func RequestLog(wrappedHandler http.Handler, lgr *loggers.AppLogger) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		lrw := NewLoggingResponseWriter(w)
		wrappedHandler.ServeHTTP(lrw, r)

		l, _ := lgr.WithReqID(r)
		start := time.Now()

		l.Info().
			Str("method", r.Method).
			Str("url", r.URL.String()).
			Str("path", r.URL.Path).
			Str("userAgent", r.UserAgent()).
			Int("respStatus", lrw.statusCode).
			Dur("elapsedMs", time.Since(start)).
			Send()
	})
}
