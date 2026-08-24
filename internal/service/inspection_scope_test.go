package service_test

import (
    "context"
    "testing"
    "time"

    "github.com/1260124186-cc/skywatch-coordinator/internal/domain"
    "github.com/1260124186-cc/skywatch-coordinator/internal/service"
    "github.com/1260124186-cc/skywatch-coordinator/internal/store"
)

func inspectionScopeCampaign(t *testing.T, c *service.Coordinator, name, code string, start time.Time) (domain.Campaign, domain.Station) {
    t.Helper()
    campaign, err := c.CreateCampaign(context.Background(), domain.CreateCampaignInput{Name: name, Target: "sky", StartsAt: start, EndsAt: start.Add(5*time.Hour)})
    if err != nil { t.Fatal(err) }
    station, err := c.AddStation(context.Background(), campaign.ID, domain.AddStationInput{Code: code, Name: name + " station", ElevationMeters: 120})
    if err != nil { t.Fatal(err) }
    return campaign, station
}

func TestInspectionReleasesCampaignScope(t *testing.T) {
    c := service.NewCoordinator(store.NewMemoryStore())
    start := time.Date(2026, 8, 25, 18, 0, 0, 0, time.UTC)
    first, _ := inspectionScopeCampaign(t, c, "North audit", "NAUD", start)
    second, _ := inspectionScopeCampaign(t, c, "South audit", "SAUD", start.Add(8*time.Hour))
    if _, err := c.InspectCampaign(context.Background(), first.ID); err != nil { t.Fatalf("first inspection: %v", err) }
    if _, err := c.InspectCampaign(context.Background(), second.ID); err != nil { t.Fatalf("independent inspection should be available: %v", err) }
}

func TestOverlapChecksReleaseCampaignScope(t *testing.T) {
    c := service.NewCoordinator(store.NewMemoryStore())
    start := time.Date(2026, 8, 25, 18, 0, 0, 0, time.UTC)
    first, firstStation := inspectionScopeCampaign(t, c, "North overlap", "NOVR", start)
    second, secondStation := inspectionScopeCampaign(t, c, "South overlap", "SOVR", start.Add(8*time.Hour))
    if _, err := c.OverlappingShifts(context.Background(), first.ID, firstStation.ID); err != nil { t.Fatalf("first overlap check: %v", err) }
    if _, err := c.OverlappingShifts(context.Background(), second.ID, secondStation.ID); err != nil { t.Fatalf("independent overlap check should be available: %v", err) }
}
