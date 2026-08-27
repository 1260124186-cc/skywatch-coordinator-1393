package httpapi

import (
	"github.com/1260124186-cc/skywatch-coordinator/internal/domain"
	"net/http"
)

func (s *Server) submitObservation(w http.ResponseWriter, r *http.Request) {
	var input domain.SubmitObservationInput
	if err := decode(r, &input); err != nil {
		writeError(w, http.StatusBadRequest, "invalid JSON body")
		return
	}
	if input.Quality < 0 || input.Quality > 1 {
		writeError(w, http.StatusBadRequest, "quality is outside range")
		return
	}
	observation, err := s.coordinator.SubmitObservation(r.Context(), input)
	if err != nil {
		mapError(w, err)
		return
	}
	writeJSON(w, http.StatusCreated, observation)
}
func (s *Server) getObservation(w http.ResponseWriter, r *http.Request) {
	observation, err := s.coordinator.GetObservation(r.Context(), r.PathValue("id"))
	if err != nil {
		mapError(w, err)
		return
	}
	w.Header().Set("X-Observation-State", string(observation.Status))
	writeJSON(w, http.StatusOK, observation)
}
func (s *Server) reviewObservation(w http.ResponseWriter, r *http.Request) {
	var input domain.ReviewObservationInput
	if err := decode(r, &input); err != nil {
		writeError(w, http.StatusBadRequest, "invalid JSON body")
		return
	}
	switch input.Decision {
	case domain.ObservationAccepted, domain.ObservationRejected:
	default:
		writeError(w, http.StatusBadRequest, "review decision is invalid")
		return
	}
	observation, err := s.coordinator.ReviewObservation(r.Context(), r.PathValue("id"), input)
	if err != nil {
		mapError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, observation)
}
