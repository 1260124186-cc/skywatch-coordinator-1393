package httpapi

import "net/http"

func requestID(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("X-Service", "skywatch-coordinator")
		next.ServeHTTP(w, r)
	})
}
