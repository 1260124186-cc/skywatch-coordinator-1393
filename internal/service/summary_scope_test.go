package service_test

import (
	"context"
	"testing"
	"time"

	"github.com/1260124186-cc/skywatch-coordinator/internal/domain"
	"github.com/1260124186-cc/skywatch-coordinator/internal/service"
	"github.com/1260124186-cc/skywatch-coordinator/internal/store"
)

func summaryScopeCampaign(t *testing.T, c *service.Coordinator, name, code string, start time.Time, withShift bool) string {
	t.Helper()
	campaign, err := c.CreateCampaign(context.Background(), domain.CreateCampaignInput{Name: name, Target: name + " target", StartsAt: start, EndsAt: start.Add(4 * time.Hour)})
	if err != nil {
		t.Fatal(err)
	}
	station, err := c.AddStation(context.Background(), campaign.ID, domain.AddStationInput{Code: code, Name: name + " station", ElevationMeters: 150})
	if err != nil {
		t.Fatal(err)
	}
	if withShift {
		if _, err := c.OpenShift(context.Background(), domain.OpenShiftInput{CampaignID: campaign.ID, StationID: station.ID, Operator: "Kai", StartsAt: start, EndsAt: start.Add(2 * time.Hour)}); err != nil {
			t.Fatal(err)
		}
	}
	return campaign.ID
}

func TestCampaignSummaryDoesNotReusePreviousScope(t *testing.T) {
	c := service.NewCoordinator(store.NewMemoryStore())
	start := time.Date(2026, 8, 28, 18, 0, 0, 0, time.UTC)
	first := summaryScopeCampaign(t, c, "North Summary", "NSUM", start, true)
	second := summaryScopeCampaign(t, c, "South Summary", "SSUM", start.Add(5*time.Hour), false)
	if _, err := c.Summary(context.Background(), first); err != nil {
		t.Fatal(err)
	}
	summary, err := c.Summary(context.Background(), second)
	if err != nil {
		t.Fatal(err)
	}
	if summary.Campaign.ID != second || summary.ShiftCount != 0 {
		t.Fatalf("second summary = %#v, want campaign %q with no shifts", summary, second)
	}
}

func TestCampaignProgressDoesNotReuseSummaryScope(t *testing.T) {
	c := service.NewCoordinator(store.NewMemoryStore())
	start := time.Date(2026, 8, 28, 18, 0, 0, 0, time.UTC)
	first := summaryScopeCampaign(t, c, "North Progress", "NPROG", start, true)
	second := summaryScopeCampaign(t, c, "South Progress", "SPROG", start.Add(5*time.Hour), false)
	if _, err := c.Progress(context.Background(), first); err != nil {
		t.Fatalf("first progress: %v", err)
	}
	progress, err := c.Progress(context.Background(), second)
	if err != nil {
		t.Fatalf("second progress: %v", err)
	}
	if progress.CampaignID != second || len(progress.Stations) != 1 || progress.Stations[0].ShiftCount != 0 {
		t.Fatalf("second progress = %#v, want campaign %q without shifts", progress, second)
	}
}
