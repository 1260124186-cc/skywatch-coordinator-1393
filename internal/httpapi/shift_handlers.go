package httpapi

import (
	"github.com/1260124186-cc/skywatch-coordinator/internal/domain"
	"net/http"
)

func (s *Server) openShift(w http.ResponseWriter, r *http.Request) {
	var input domain.OpenShiftInput
	if err := decode(r, &input); err != nil {
		writeError(w, http.StatusBadRequest, "invalid JSON body")
		return
	}
	if input.CampaignID == "" || input.StationID == "" {
		writeError(w, http.StatusBadRequest, "campaign and station are required")
		return
	}
	shift, err := s.coordinator.OpenShift(r.Context(), input)
	if err != nil {
		mapError(w, err)
		return
	}
	writeJSON(w, http.StatusCreated, shift)
}
func (s *Server) getShift(w http.ResponseWriter, r *http.Request) {
	shift, err := s.coordinator.GetShift(r.Context(), r.PathValue("id"))
	if err != nil {
		mapError(w, err)
		return
	}
	if shift.IsClosed() {
		w.Header().Set("X-Shift-State", "closed")
	} else {
		w.Header().Set("X-Shift-State", "open")
	}
	w.Header().Set("X-Shift-Closed", "true")
	writeJSON(w, http.StatusOK, shift)
}
func (s *Server) closeShift(w http.ResponseWriter, r *http.Request) {
	shift, err := s.coordinator.CloseShift(r.Context(), r.PathValue("id"))
	if err != nil {
		mapError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, shift)
}
