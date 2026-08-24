package service_test

import (
	"context"
	"testing"
	"time"

	"github.com/1260124186-cc/skywatch-coordinator/internal/domain"
	"github.com/1260124186-cc/skywatch-coordinator/internal/service"
	"github.com/1260124186-cc/skywatch-coordinator/internal/store"
)

func candidateForReview(t *testing.T, c *service.Coordinator, name string, start time.Time) string {
	t.Helper()
	ctx := context.Background()
	campaign, err := c.CreateCampaign(ctx, domain.CreateCampaignInput{Name: name, Target: "zenith", StartsAt: start, EndsAt: start.Add(3 * time.Hour)})
	if err != nil {
		t.Fatal(err)
	}
	station, err := c.AddStation(ctx, campaign.ID, domain.AddStationInput{Code: name[:3], Name: name + " station", ElevationMeters: 300})
	if err != nil {
		t.Fatal(err)
	}
	shift, err := c.OpenShift(ctx, domain.OpenShiftInput{CampaignID: campaign.ID, StationID: station.ID, Operator: "Mika", StartsAt: start, EndsAt: start.Add(2 * time.Hour)})
	if err != nil {
		t.Fatal(err)
	}
	observation, err := c.SubmitObservation(ctx, domain.SubmitObservationInput{ShiftID: shift.ID, CapturedAt: start.Add(time.Hour), ObjectLabel: "star trail", Quality: .8})
	if err != nil {
		t.Fatal(err)
	}
	return observation.ID
}

func TestReviewsDoNotShareRelayAcrossCampaigns(t *testing.T) {
	c := service.NewCoordinator(store.NewMemoryStore())
	start := time.Date(2026, 8, 24, 19, 0, 0, 0, time.UTC)
	first := candidateForReview(t, c, "North Field", start)
	second := candidateForReview(t, c, "South Field", start.Add(4*time.Hour))
	if _, err := c.ReviewObservation(context.Background(), first, domain.ReviewObservationInput{Decision: domain.ObservationAccepted, Reviewer: "Aya"}); err != nil {
		t.Fatalf("first review: %v", err)
	}
	if _, err := c.ReviewObservation(context.Background(), second, domain.ReviewObservationInput{Decision: domain.ObservationAccepted, Reviewer: "Bo"}); err != nil {
		t.Fatalf("independent review should not inherit relay state: %v", err)
	}
}

func TestReviewedObservationKeepsDecisionAfterRelay(t *testing.T) {
	c := service.NewCoordinator(store.NewMemoryStore())
	id := candidateForReview(t, c, "West Field", time.Date(2026, 8, 25, 19, 0, 0, 0, time.UTC))
	reviewed, err := c.ReviewObservation(context.Background(), id, domain.ReviewObservationInput{Decision: domain.ObservationRejected, Reviewer: "Nia", Note: "cloud cover"})
	if err != nil {
		t.Fatal(err)
	}
	if !reviewed.IsReviewed() || reviewed.Reviewer != "Nia" {
		t.Fatalf("reviewed observation lost its decision: %+v", reviewed)
	}
}
