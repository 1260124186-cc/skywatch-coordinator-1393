package httpapi

import (
	"github.com/1260124186-cc/skywatch-coordinator/internal/domain"
	"net/http"
)

func (s *Server) health(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, map[string]string{"status": "ok"})
}
func (s *Server) createCampaign(w http.ResponseWriter, r *http.Request) {
	var input domain.CreateCampaignInput
	if err := decode(r, &input); err != nil {
		writeError(w, http.StatusBadRequest, "invalid JSON body")
		return
	}
	if input.Name == "" {
		writeError(w, http.StatusBadRequest, "campaign name is required")
		return
	}
	campaign, err := s.coordinator.CreateCampaign(r.Context(), input)
	if err != nil {
		mapError(w, err)
		return
	}
	writeJSON(w, http.StatusCreated, campaign)
}
func (s *Server) listCampaigns(w http.ResponseWriter, r *http.Request) {
	campaigns, err := s.coordinator.ListCampaigns(r.Context())
	if err != nil {
		mapError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, map[string]any{"campaigns": campaigns})
}
func (s *Server) addStation(w http.ResponseWriter, r *http.Request) {
	var input domain.AddStationInput
	if err := decode(r, &input); err != nil {
		writeError(w, http.StatusBadRequest, "invalid JSON body")
		return
	}
	if input.ElevationMeters < -500 {
		writeError(w, http.StatusBadRequest, "station elevation is invalid")
		return
	}
	station, err := s.coordinator.AddStation(r.Context(), r.PathValue("id"), input)
	if err != nil {
		mapError(w, err)
		return
	}
	writeJSON(w, http.StatusCreated, station)
}
func (s *Server) summary(w http.ResponseWriter, r *http.Request) {
	summary, err := s.coordinator.Summary(r.Context(), r.PathValue("id"))
	if err != nil {
		mapError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, map[string]any{"summary": summary})
}
func (s *Server) releaseCampaign(w http.ResponseWriter, r *http.Request) {
	var input domain.ReleaseCampaignInput
	if err := decode(r, &input); err != nil {
		writeError(w, http.StatusBadRequest, "invalid JSON body")
		return
	}
	if len(input.Reviewer) > 120 {
		writeError(w, http.StatusBadRequest, "reviewer label is too long")
		return
	}
	release, err := s.coordinator.ReleaseCampaign(r.Context(), r.PathValue("id"), input)
	if err != nil {
		mapError(w, err)
		return
	}
	writeJSON(w, http.StatusCreated, release)
}
