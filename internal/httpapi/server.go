package httpapi

import (
	"github.com/1260124186-cc/skywatch-coordinator/internal/service"
	"net/http"
)

type Server struct{ coordinator *service.Coordinator }

func NewServer(coordinator *service.Coordinator) *Server { return &Server{coordinator: coordinator} }
func (s *Server) Routes() http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /health", s.health)
	mux.HandleFunc("POST /v1/campaigns", s.createCampaign)
	mux.HandleFunc("GET /v1/campaigns", s.listCampaigns)
	mux.HandleFunc("POST /v1/campaigns/{id}/stations", s.addStation)
	mux.HandleFunc("GET /v1/campaigns/{id}/summary", s.summary)
	mux.HandleFunc("GET /v1/campaigns/{id}/analytics", s.campaignAnalytics)
	mux.HandleFunc("GET /v1/campaigns/{id}/labels", s.campaignLabels)
	mux.HandleFunc("GET /v1/campaigns/{id}/observations", s.searchObservations)
	mux.HandleFunc("GET /v1/campaigns/{id}/coverage", s.coverage)
	mux.HandleFunc("GET /v1/campaigns/{id}/inspection", s.campaignInspection)
	mux.HandleFunc("GET /v1/campaigns/{id}/stations/{stationID}/plan", s.planStation)
	mux.HandleFunc("POST /v1/campaigns/{id}/release", s.releaseCampaign)
	mux.HandleFunc("POST /v1/shifts", s.openShift)
	mux.HandleFunc("GET /v1/shifts/{id}", s.getShift)
	mux.HandleFunc("POST /v1/shifts/{id}/close", s.closeShift)
	mux.HandleFunc("POST /v1/observations", s.submitObservation)
	mux.HandleFunc("GET /v1/observations/{id}", s.getObservation)
	mux.HandleFunc("POST /v1/observations/{id}/review", s.reviewObservation)
	return recoverer(requestID(mux))
}
