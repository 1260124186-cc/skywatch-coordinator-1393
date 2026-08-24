package domain

import "time"

type CampaignStatus string

const (
	CampaignActive   CampaignStatus = "active"
	CampaignReleased CampaignStatus = "released"
	CampaignArchived CampaignStatus = "archived"
)

type Campaign struct {
	ID         string         `json:"id"`
	Name       string         `json:"name"`
	Target     string         `json:"target"`
	StartsAt   time.Time      `json:"startsAt"`
	EndsAt     time.Time      `json:"endsAt"`
	Status     CampaignStatus `json:"status"`
	CreatedAt  time.Time      `json:"createdAt"`
	ReleasedAt *time.Time     `json:"releasedAt,omitempty"`
}

func (c Campaign) Contains(at time.Time) bool { return !at.Before(c.StartsAt) && !at.After(c.EndsAt) }
func (c Campaign) CanAcceptShift() bool       { return c.Status == CampaignActive }
func (c Campaign) Duration() time.Duration    { return c.EndsAt.Sub(c.StartsAt) }
func (c Campaign) IsReleased() bool           { return c.Status == CampaignReleased }
