package middleware

import (
	"log"
	"net/http"
	"time"
)

type statusRecorder struct {
	http.ResponseWriter
	status int
}

func (r *statusRecorder) WriteHeader(statusCode int) {
	r.status = statusCode
	r.ResponseWriter.WriteHeader(statusCode)
}

func (r *statusRecorder) Status() int {
	return r.status
}

func Logger(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		startTime := time.Now()
		recorder := &statusRecorder{ResponseWriter: w, status: http.StatusOK}
		next.ServeHTTP(recorder, r)
		endTime := time.Now()
		duration := endTime.Sub(startTime)
		statusCode := recorder.Status()
		endpoint := r.URL.Path
		log.Printf("[%s] %s - %v - %d\n", r.Method, endpoint, duration, statusCode)
	})
}
