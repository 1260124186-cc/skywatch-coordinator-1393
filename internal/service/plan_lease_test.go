package service_test

import (
	"context"
	"testing"
	"time"

	"github.com/1260124186-cc/skywatch-coordinator/internal/domain"
	"github.com/1260124186-cc/skywatch-coordinator/internal/service"
	"github.com/1260124186-cc/skywatch-coordinator/internal/store"
)

func planLeaseCampaign(t *testing.T, c *service.Coordinator, name, code string, start time.Time) (string, string) {
	t.Helper()
	campaign, err := c.CreateCampaign(context.Background(), domain.CreateCampaignInput{Name: name, Target: "night sky", StartsAt: start, EndsAt: start.Add(6 * time.Hour)})
	if err != nil {
		t.Fatal(err)
	}
	station, err := c.AddStation(context.Background(), campaign.ID, domain.AddStationInput{Code: code, Name: name + " station", ElevationMeters: 420})
	if err != nil {
		t.Fatal(err)
	}
	return campaign.ID, station.ID
}

func TestStationPlansReleaseCompletedLease(t *testing.T) {
	c := service.NewCoordinator(store.NewMemoryStore())
	start := time.Date(2026, 8, 25, 18, 0, 0, 0, time.UTC)
	firstCampaign, firstStation := planLeaseCampaign(t, c, "North Plan", "NPLAN", start)
	secondCampaign, secondStation := planLeaseCampaign(t, c, "South Plan", "SPLAN", start.Add(8*time.Hour))
	if _, err := c.PlanForStation(context.Background(), firstCampaign, firstStation, 3); err != nil {
		t.Fatalf("first plan: %v", err)
	}
	plan, err := c.PlanForStation(context.Background(), secondCampaign, secondStation, 4)
	if err != nil {
		t.Fatalf("independent second plan should be available: %v", err)
	}
	if plan.CampaignID != secondCampaign || plan.StationID != secondStation || len(plan.Windows) != 4 {
		t.Fatalf("second plan = %#v", plan)
	}
}

func TestAvailableWindowsReleaseCompletedLease(t *testing.T) {
	c := service.NewCoordinator(store.NewMemoryStore())
	start := time.Date(2026, 8, 25, 18, 0, 0, 0, time.UTC)
	firstCampaign, firstStation := planLeaseCampaign(t, c, "North Window", "NWIN", start)
	secondCampaign, secondStation := planLeaseCampaign(t, c, "South Window", "SWIN", start.Add(8*time.Hour))
	if _, err := c.AvailableWindows(context.Background(), firstCampaign, firstStation); err != nil {
		t.Fatalf("first windows: %v", err)
	}
	windows, err := c.AvailableWindows(context.Background(), secondCampaign, secondStation)
	if err != nil {
		t.Fatalf("independent second windows should be available: %v", err)
	}
	if len(windows) != 1 || windows[0].Duration() != 6*time.Hour {
		t.Fatalf("second windows = %#v", windows)
	}
}
