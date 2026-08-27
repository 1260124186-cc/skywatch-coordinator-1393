package domain

import "time"

type ObservationStatus string

const (
	ObservationCandidate ObservationStatus = "candidate"
	ObservationAccepted  ObservationStatus = "accepted"
	ObservationRejected  ObservationStatus = "rejected"
)

type Observation struct {
	ID          string            `json:"id"`
	CampaignID  string            `json:"campaignId"`
	ShiftID     string            `json:"shiftId"`
	CapturedAt  time.Time         `json:"capturedAt"`
	ObjectLabel string            `json:"objectLabel"`
	Quality     float64           `json:"quality"`
	Status      ObservationStatus `json:"status"`
	Reviewer    string            `json:"reviewer,omitempty"`
	ReviewNote  string            `json:"reviewNote,omitempty"`
}

func (o Observation) IsCandidate() bool { return o.Status == ObservationCandidate }
func (o Observation) IsAccepted() bool  { return o.Status == ObservationAccepted }
func (o Observation) IsReviewed() bool  { return o.Status != ObservationCandidate }
