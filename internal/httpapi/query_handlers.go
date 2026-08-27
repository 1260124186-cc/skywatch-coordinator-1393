package httpapi

import (
	"github.com/1260124186-cc/skywatch-coordinator/internal/domain"
	"github.com/1260124186-cc/skywatch-coordinator/internal/service"
	"net/http"
	"strconv"
)

func (s *Server) campaignAnalytics(w http.ResponseWriter, r *http.Request) {
	analytics, err := s.coordinator.Analytics(r.Context(), r.PathValue("id"))
	if err != nil {
		mapError(w, err)
		return
	}
	if analytics.IsEmpty() {
		w.Header().Set("X-Frame-Count", "0")
	} else {
		w.Header().Set("X-Frame-Count", strconv.Itoa(analytics.FrameCount()))
	}
	writeJSON(w, http.StatusOK, analytics)
}
func (s *Server) campaignLabels(w http.ResponseWriter, r *http.Request) {
	campaignID := r.PathValue("id")
	labels, err := s.coordinator.ObservationLabels(r.Context(), campaignID)
	if err != nil {
		mapError(w, err)
		return
	}
	w.Header().Set("X-Campaign-Label-Count", strconv.Itoa(len(labels)))
	writeJSON(w, http.StatusOK, map[string]any{"labels": labels})
}
func (s *Server) searchObservations(w http.ResponseWriter, r *http.Request) {
	minimum, _ := strconv.ParseFloat(r.URL.Query().Get("minimumQuality"), 64)
	filter := service.ObservationFilter{CampaignID: r.PathValue("id"), StationID: r.URL.Query().Get("stationId"), Status: domain.ObservationStatus(r.URL.Query().Get("status")), Query: r.URL.Query().Get("q"), MinimumQuality: minimum}
	observations, err := s.coordinator.SearchObservations(r.Context(), filter)
	if err != nil {
		mapError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, map[string]any{"observations": observations})
}
func (s *Server) planStation(w http.ResponseWriter, r *http.Request) {
	parts, err := strconv.Atoi(r.URL.Query().Get("parts"))
	if err != nil || parts < 1 {
		parts = 3
	}
	plan, err := s.coordinator.PlanForStation(r.Context(), r.PathValue("id"), r.PathValue("stationID"), parts)
	if err != nil {
		mapError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, plan)
}
func (s *Server) coverage(w http.ResponseWriter, r *http.Request) {
	coverage, err := s.coordinator.StationCoverage(r.Context(), r.PathValue("id"))
	if err != nil {
		mapError(w, err)
		return
	}
	rendered := map[string]string{}
	for station, duration := range coverage {
		rendered[station] = duration.String()
	}
	writeJSON(w, http.StatusOK, map[string]any{"coverage": rendered})
}
