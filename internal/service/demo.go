package service

import (
	"github.com/1260124186-cc/skywatch-coordinator/internal/domain"
	"time"
)

func LoadDemo(c *Coordinator) error {
	start := time.Date(2026, 8, 24, 18, 0, 0, 0, time.UTC)
	end := start.Add(10 * time.Hour)
	campaign := domain.Campaign{ID: "cmp-north", Name: "Aurora Sweep", Target: "northern horizon", StartsAt: start, EndsAt: end, Status: domain.CampaignActive, CreatedAt: start}
	if err := c.repo.CreateCampaign(campaign); err != nil {
		return err
	}
	station := domain.Station{ID: "stn-north", CampaignID: campaign.ID, Code: "NORTH-01", Name: "Ridge Station", ElevationMeters: 840, Active: true}
	if err := c.repo.AddStation(station); err != nil {
		return err
	}
	shift := domain.Shift{ID: "shift-north", CampaignID: campaign.ID, StationID: station.ID, Operator: "Mika", StartsAt: start, EndsAt: end, Status: domain.ShiftOpen}
	if err := c.repo.CreateShift(shift); err != nil {
		return err
	}
	observation := domain.Observation{ID: "obs-north", CampaignID: campaign.ID, ShiftID: shift.ID, CapturedAt: start.Add(time.Hour), ObjectLabel: "aurora arc", Quality: 0.82, Status: domain.ObservationCandidate}
	return c.repo.CreateObservation(observation)
}
