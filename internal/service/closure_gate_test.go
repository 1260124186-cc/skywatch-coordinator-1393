package service_test

import (
	"context"
	"testing"
	"time"

	"github.com/1260124186-cc/skywatch-coordinator/internal/domain"
	"github.com/1260124186-cc/skywatch-coordinator/internal/service"
	"github.com/1260124186-cc/skywatch-coordinator/internal/store"
)

func openShiftForClosure(t *testing.T, c *service.Coordinator, name string, start time.Time) string {
	t.Helper()
	ctx := context.Background()
	campaign, err := c.CreateCampaign(ctx, domain.CreateCampaignInput{Name: name, Target: "horizon", StartsAt: start, EndsAt: start.Add(3 * time.Hour)})
	if err != nil {
		t.Fatal(err)
	}
	station, err := c.AddStation(ctx, campaign.ID, domain.AddStationInput{Code: name[:3], Name: name + " station", ElevationMeters: 400})
	if err != nil {
		t.Fatal(err)
	}
	shift, err := c.OpenShift(ctx, domain.OpenShiftInput{CampaignID: campaign.ID, StationID: station.ID, Operator: "Mika", StartsAt: start, EndsAt: start.Add(2 * time.Hour)})
	if err != nil {
		t.Fatal(err)
	}
	return shift.ID
}

func TestShiftClosuresStayIndependentAcrossCampaigns(t *testing.T) {
	c := service.NewCoordinator(store.NewMemoryStore())
	start := time.Date(2026, 8, 24, 20, 0, 0, 0, time.UTC)
	first := openShiftForClosure(t, c, "North Field", start)
	second := openShiftForClosure(t, c, "South Field", start.Add(4*time.Hour))
	if _, err := c.CloseShift(context.Background(), first); err != nil {
		t.Fatalf("first close: %v", err)
	}
	if _, err := c.CloseShift(context.Background(), second); err != nil {
		t.Fatalf("independent close should not inherit closure state: %v", err)
	}
}

func TestClosedShiftRetainsCloseTime(t *testing.T) {
	c := service.NewCoordinator(store.NewMemoryStore())
	id := openShiftForClosure(t, c, "West Field", time.Date(2026, 8, 25, 20, 0, 0, 0, time.UTC))
	shift, err := c.CloseShift(context.Background(), id)
	if err != nil {
		t.Fatal(err)
	}
	if !shift.IsClosed() || shift.ClosedAt == nil {
		t.Fatalf("close result lost state: %+v", shift)
	}
}
