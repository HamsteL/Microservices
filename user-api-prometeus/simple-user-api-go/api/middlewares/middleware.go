package middlewares

import (
	"github.com/prometheus/client_golang/prometheus"
	"net/http"
	"strconv"
)

func PrometheusMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		timer := prometheus.NewTimer(httpDuration.WithLabelValues(r.RequestURI))

		lrw := NewResponseWriter(w)
		next.ServeHTTP(lrw, r)

		totalRequests.WithLabelValues(r.RequestURI).Inc()
		responseStatus.WithLabelValues(r.RequestURI, strconv.Itoa(lrw.statusCode)).Inc()
		timer.ObserveDuration()
	})
}

type responseWriter struct {
	http.ResponseWriter
	statusCode int
}

func NewResponseWriter(w http.ResponseWriter) *responseWriter {
	return &responseWriter{w, http.StatusOK}
}

func (rw *responseWriter) WriteHeader(code int) {
	rw.statusCode = code
	rw.ResponseWriter.WriteHeader(code)
}
