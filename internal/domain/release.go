package domain

import "time"

type Release struct {
	ID            string    `json:"id"`
	CampaignID    string    `json:"campaignId"`
	AcceptedCount int       `json:"acceptedCount"`
	RejectedCount int       `json:"rejectedCount"`
	ReleasedBy    string    `json:"releasedBy"`
	ReleasedAt    time.Time `json:"releasedAt"`
}

func (r Release) TotalReviewed() int      { return r.AcceptedCount + r.RejectedCount }
func (r Release) HasAcceptedFrames() bool { return r.AcceptedCount > 0 }
