package middlewares

import (
	"net/http"
	"time"

	"github.com/go-chi/chi/v5/middleware"
	"github.com/happYness-Project/taskManagementGolang/pkg/loggers"
)

func Logger(logger *loggers.AppLogger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		fn := func(rw http.ResponseWriter, r *http.Request) {
			ww := middleware.NewWrapResponseWriter(rw, r.ProtoMajor)
			start := time.Now()
			defer func() {
				logger.Info().
					// Str("request-id", GetReqID(r.Context())).
					Int("status", ww.Status()).
					Int("bytes", ww.BytesWritten()).
					Str("method", r.Method).
					Str("path", r.URL.Path).
					Str("query", r.URL.RawQuery).
					Str("ip", r.RemoteAddr).
					// Str("trace.id", trace.SpanFromContext(r.Context()).SpanContext().TraceID().String()).
					Str("user-agent", r.UserAgent()).
					Dur("latency", time.Since(start)).
					Msg("request completed")
			}()

			next.ServeHTTP(ww, r)
		}
		return http.HandlerFunc(fn)
	}
}

// func RequestLog(wrappedHandler http.Handler, lgr *loggers.AppLogger) http.Handler {
// 	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

// 		lrw := NewLoggingResponseWriter(w)
// 		wrappedHandler.ServeHTTP(lrw, r)

// 		l, _ := lgr.WithReqID(r)
// 		start := time.Now()

// 		l.Info().
// 			Str("method", r.Method).
// 			Str("url", r.URL.String()).
// 			Str("path", r.URL.Path).
// 			Str("userAgent", r.UserAgent()).
// 			Int("respStatus", lrw.statusCode).
// 			Dur("elapsedMs", time.Since(start)).
// 			Send()
// 	})
// }
