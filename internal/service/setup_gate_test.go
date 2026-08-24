package service_test

import (
    "context"
    "testing"
    "time"

    "github.com/1260124186-cc/skywatch-coordinator/internal/domain"
    "github.com/1260124186-cc/skywatch-coordinator/internal/service"
    "github.com/1260124186-cc/skywatch-coordinator/internal/store"
)

func setupGateCampaign(t *testing.T, c *service.Coordinator, name string, start time.Time) domain.Campaign {
    t.Helper()
    campaign, err := c.CreateCampaign(context.Background(), domain.CreateCampaignInput{Name: name, Target: "sky", StartsAt: start, EndsAt: start.Add(5*time.Hour)})
    if err != nil { t.Fatal(err) }
    return campaign
}

func TestStationsReleaseCompletedSetupGate(t *testing.T) {
    c := service.NewCoordinator(store.NewMemoryStore())
    start := time.Date(2026, 8, 25, 18, 0, 0, 0, time.UTC)
    first := setupGateCampaign(t, c, "North setup", start)
    second := setupGateCampaign(t, c, "South setup", start.Add(8*time.Hour))
    if _, err := c.AddStation(context.Background(), first.ID, domain.AddStationInput{Code: "NSET", Name: "North", ElevationMeters: 120}); err != nil { t.Fatalf("first setup: %v", err) }
    if _, err := c.AddStation(context.Background(), second.ID, domain.AddStationInput{Code: "SSET", Name: "South", ElevationMeters: 220}); err != nil { t.Fatalf("independent station setup should not remain busy: %v", err) }
}

func TestShiftsReleaseCompletedSetupGate(t *testing.T) {
    c := service.NewCoordinator(store.NewMemoryStore())
    start := time.Date(2026, 8, 25, 18, 0, 0, 0, time.UTC)
    first := setupGateCampaign(t, c, "First shift", start)
    second := setupGateCampaign(t, c, "Second shift", start.Add(8*time.Hour))
    if err := c.Repository().AddStation(domain.Station{ID: "north-gate", CampaignID: first.ID, Code: "N-GATE", Name: "North", ElevationMeters: 120, Active: true}); err != nil { t.Fatal(err) }
    if err := c.Repository().AddStation(domain.Station{ID: "south-gate", CampaignID: second.ID, Code: "S-GATE", Name: "South", ElevationMeters: 220, Active: true}); err != nil { t.Fatal(err) }
    if _, err := c.OpenShift(context.Background(), domain.OpenShiftInput{CampaignID: first.ID, StationID: "north-gate", Operator: "Nia", StartsAt: start, EndsAt: start.Add(time.Hour)}); err != nil { t.Fatalf("first shift: %v", err) }
    if _, err := c.OpenShift(context.Background(), domain.OpenShiftInput{CampaignID: second.ID, StationID: "south-gate", Operator: "Sol", StartsAt: start.Add(8*time.Hour), EndsAt: start.Add(9*time.Hour)}); err != nil { t.Fatalf("independent shift should not remain busy: %v", err) }
}
