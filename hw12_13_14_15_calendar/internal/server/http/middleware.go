package internalhttp

import (
	"context"
	"net/http"
	"time"
)

const timeLayout = "02/Jan/2006:15:04:05 -0700"

type statusRecorder struct {
	http.ResponseWriter
	status int
}

func (rec *statusRecorder) WriteHeader(code int) {
	rec.status = code
	rec.ResponseWriter.WriteHeader(code)
}

func loggingMiddleware(ctx context.Context, logger Logger, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		startTime := time.Now()
		rec := statusRecorder{w, 200}
		next.ServeHTTP(&rec, r)
		elapsed := time.Since(startTime)

		logger.Info(ctx, "http request",
			"ip", r.RemoteAddr,
			"startTime", startTime.Format(timeLayout),
			"method", r.Method,
			"path", r.URL.Path,
			"version", r.Proto,
			"statusCode", rec.status,
			"latency", elapsed.Milliseconds(),
			"userAgent", r.UserAgent(),
		)
	})
}
