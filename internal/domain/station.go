package domain

type Station struct {
	ID              string `json:"id"`
	CampaignID      string `json:"campaignId"`
	Code            string `json:"code"`
	Name            string `json:"name"`
	ElevationMeters int    `json:"elevationMeters"`
	Active          bool   `json:"active"`
}

func (s Station) IsUsable() bool                   { return s.Active && s.Code != "" }
func (s Station) BelongsTo(campaignID string) bool { return s.CampaignID == campaignID }
func (s Station) LocationLabel() string            { return s.Code + " · " + s.Name }
