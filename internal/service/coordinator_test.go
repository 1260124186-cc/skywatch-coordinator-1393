package service_test

import (
	"context"
	"github.com/1260124186-cc/skywatch-coordinator/internal/domain"
	"github.com/1260124186-cc/skywatch-coordinator/internal/service"
	"github.com/1260124186-cc/skywatch-coordinator/internal/store"
	"testing"
	"time"
)

func TestCampaignLifecycle(t *testing.T) {
	ctx := context.Background()
	c := service.NewCoordinator(store.NewMemoryStore())
	start := time.Date(2026, 8, 24, 18, 0, 0, 0, time.UTC)
	campaign, err := c.CreateCampaign(ctx, domain.CreateCampaignInput{Name: "Perseid Field", Target: "Perseus", StartsAt: start, EndsAt: start.Add(4 * time.Hour)})
	if err != nil {
		t.Fatal(err)
	}
	station, err := c.AddStation(ctx, campaign.ID, domain.AddStationInput{Code: "RIDGE", Name: "Ridge", ElevationMeters: 500})
	if err != nil {
		t.Fatal(err)
	}
	shift, err := c.OpenShift(ctx, domain.OpenShiftInput{CampaignID: campaign.ID, StationID: station.ID, Operator: "Lin", StartsAt: start, EndsAt: start.Add(3 * time.Hour)})
	if err != nil {
		t.Fatal(err)
	}
	observation, err := c.SubmitObservation(ctx, domain.SubmitObservationInput{ShiftID: shift.ID, CapturedAt: start.Add(time.Hour), ObjectLabel: "meteor", Quality: .8})
	if err != nil {
		t.Fatal(err)
	}
	if _, err = c.ReviewObservation(ctx, observation.ID, domain.ReviewObservationInput{Decision: domain.ObservationAccepted, Reviewer: "Aya"}); err != nil {
		t.Fatal(err)
	}
	if _, err = c.CloseShift(ctx, shift.ID); err != nil {
		t.Fatal(err)
	}
	release, err := c.ReleaseCampaign(ctx, campaign.ID, domain.ReleaseCampaignInput{Reviewer: "Aya"})
	if err != nil {
		t.Fatal(err)
	}
	if !release.HasAcceptedFrames() {
		t.Fatal("expected release with accepted frame")
	}
}
func TestFrameOutsideShiftFails(t *testing.T) {
	ctx := context.Background()
	c := service.NewCoordinator(store.NewMemoryStore())
	start := time.Date(2026, 8, 24, 18, 0, 0, 0, time.UTC)
	campaign, _ := c.CreateCampaign(ctx, domain.CreateCampaignInput{Name: "N", Target: "T", StartsAt: start, EndsAt: start.Add(4 * time.Hour)})
	station, _ := c.AddStation(ctx, campaign.ID, domain.AddStationInput{Code: "AAA", Name: "A", ElevationMeters: 1})
	shift, _ := c.OpenShift(ctx, domain.OpenShiftInput{CampaignID: campaign.ID, StationID: station.ID, Operator: "O", StartsAt: start, EndsAt: start.Add(time.Hour)})
	_, err := c.SubmitObservation(ctx, domain.SubmitObservationInput{ShiftID: shift.ID, CapturedAt: start.Add(2 * time.Hour), ObjectLabel: "x", Quality: .4})
	if err == nil {
		t.Fatal("expected conflict")
	}
}
