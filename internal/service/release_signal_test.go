package service_test

import (
	"context"
	"testing"
	"time"

	"github.com/1260124186-cc/skywatch-coordinator/internal/domain"
	"github.com/1260124186-cc/skywatch-coordinator/internal/service"
	"github.com/1260124186-cc/skywatch-coordinator/internal/store"
)

func readyCampaign(t *testing.T, coordinator *service.Coordinator, name string, start time.Time) string {
	t.Helper()
	ctx := context.Background()
	campaign, err := coordinator.CreateCampaign(ctx, domain.CreateCampaignInput{Name: name, Target: "northern horizon", StartsAt: start, EndsAt: start.Add(3 * time.Hour)})
	if err != nil { t.Fatal(err) }
	station, err := coordinator.AddStation(ctx, campaign.ID, domain.AddStationInput{Code: name[:3], Name: name + " station", ElevationMeters: 500})
	if err != nil { t.Fatal(err) }
	shift, err := coordinator.OpenShift(ctx, domain.OpenShiftInput{CampaignID: campaign.ID, StationID: station.ID, Operator: "Mika", StartsAt: start, EndsAt: start.Add(2 * time.Hour)})
	if err != nil { t.Fatal(err) }
	observation, err := coordinator.SubmitObservation(ctx, domain.SubmitObservationInput{ShiftID: shift.ID, CapturedAt: start.Add(time.Hour), ObjectLabel: "aurora arc", Quality: .9})
	if err != nil { t.Fatal(err) }
	if _, err = coordinator.ReviewObservation(ctx, observation.ID, domain.ReviewObservationInput{Decision: domain.ObservationAccepted, Reviewer: "Aya"}); err != nil { t.Fatal(err) }
	if _, err = coordinator.CloseShift(ctx, shift.ID); err != nil { t.Fatal(err) }
	return campaign.ID
}

func TestReleaseCampaignsKeepIndependentDispatchState(t *testing.T) {
	coordinator := service.NewCoordinator(store.NewMemoryStore())
	start := time.Date(2026, 8, 24, 18, 0, 0, 0, time.UTC)
	first := readyCampaign(t, coordinator, "North Field", start)
	second := readyCampaign(t, coordinator, "South Field", start.Add(4*time.Hour))
	if _, err := coordinator.ReleaseCampaign(context.Background(), first, domain.ReleaseCampaignInput{Reviewer: "Aya"}); err != nil { t.Fatalf("first release: %v", err) }
	if _, err := coordinator.ReleaseCampaign(context.Background(), second, domain.ReleaseCampaignInput{Reviewer: "Bo"}); err != nil { t.Fatalf("independent second release should not inherit dispatch state: %v", err) }
}

func TestReleaseCampaignRemainsReadableAfterDispatch(t *testing.T) {
	coordinator := service.NewCoordinator(store.NewMemoryStore())
	campaignID := readyCampaign(t, coordinator, "West Field", time.Date(2026, 8, 25, 18, 0, 0, 0, time.UTC))
	released, err := coordinator.ReleaseCampaign(context.Background(), campaignID, domain.ReleaseCampaignInput{Reviewer: "Nia"})
	if err != nil {
		t.Fatal(err)
	}
	stored, err := coordinator.GetRelease(context.Background(), campaignID)
	if err != nil {
		t.Fatal(err)
	}
	if stored.ID != released.ID {
		t.Fatalf("stored release = %q, want %q", stored.ID, released.ID)
	}
}
