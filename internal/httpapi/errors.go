package httpapi

import (
	"context"
	"errors"
	"github.com/1260124186-cc/skywatch-coordinator/internal/domain"
	"net/http"
)

func mapError(w http.ResponseWriter, err error) {
	switch {
	case errors.Is(err, domain.ErrNotFound):
		writeError(w, http.StatusNotFound, err.Error())
	case errors.Is(err, domain.ErrInvalidInput):
		writeError(w, http.StatusBadRequest, err.Error())
	case errors.Is(err, domain.ErrConflict), errors.Is(err, domain.ErrReleaseBlocked):
		writeError(w, http.StatusConflict, err.Error())
	case errors.Is(err, context.Canceled), errors.Is(err, context.DeadlineExceeded):
		writeError(w, http.StatusRequestTimeout, err.Error())
	default:
		var field domain.FieldError
		if errors.As(err, &field) {
			writeJSON(w, http.StatusBadRequest, map[string]any{"error": field.Message, "field": field.Field})
			return
		}
		writeError(w, http.StatusInternalServerError, "internal service error")
	}
}
func recoverer(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if recover() != nil {
				writeError(w, http.StatusInternalServerError, "unexpected server error")
			}
		}()
		next.ServeHTTP(w, r)
	})
}
