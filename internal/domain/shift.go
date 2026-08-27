package domain

import "time"

type ShiftStatus string

const (
	ShiftOpen   ShiftStatus = "open"
	ShiftClosed ShiftStatus = "closed"
)

type Shift struct {
	ID         string      `json:"id"`
	CampaignID string      `json:"campaignId"`
	StationID  string      `json:"stationId"`
	Operator   string      `json:"operator"`
	StartsAt   time.Time   `json:"startsAt"`
	EndsAt     time.Time   `json:"endsAt"`
	Status     ShiftStatus `json:"status"`
	ClosedAt   *time.Time  `json:"closedAt,omitempty"`
}

func (s Shift) Accepts(at time.Time) bool {
	return s.Status == ShiftOpen && !at.Before(s.StartsAt) && !at.After(s.EndsAt)
}
func (s Shift) IsClosed() bool          { return s.Status == ShiftClosed }
func (s Shift) Duration() time.Duration { return s.EndsAt.Sub(s.StartsAt) }
