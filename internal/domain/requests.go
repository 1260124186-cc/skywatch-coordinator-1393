package domain

import "time"

type CreateCampaignInput struct {
	Name     string    `json:"name"`
	Target   string    `json:"target"`
	StartsAt time.Time `json:"startsAt"`
	EndsAt   time.Time `json:"endsAt"`
}
type AddStationInput struct {
	Code            string `json:"code"`
	Name            string `json:"name"`
	ElevationMeters int    `json:"elevationMeters"`
}
type OpenShiftInput struct {
	CampaignID string    `json:"campaignId"`
	StationID  string    `json:"stationId"`
	Operator   string    `json:"operator"`
	StartsAt   time.Time `json:"startsAt"`
	EndsAt     time.Time `json:"endsAt"`
}
type SubmitObservationInput struct {
	ShiftID     string    `json:"shiftId"`
	CapturedAt  time.Time `json:"capturedAt"`
	ObjectLabel string    `json:"objectLabel"`
	Quality     float64   `json:"quality"`
}
type ReviewObservationInput struct {
	Decision ObservationStatus `json:"decision"`
	Reviewer string            `json:"reviewer"`
	Note     string            `json:"note"`
}
type ReleaseCampaignInput struct {
	Reviewer string `json:"reviewer"`
}
