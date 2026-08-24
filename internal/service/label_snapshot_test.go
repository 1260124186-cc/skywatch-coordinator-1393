package service_test

import (
	"context"
	"github.com/1260124186-cc/skywatch-coordinator/internal/domain"
	"github.com/1260124186-cc/skywatch-coordinator/internal/service"
	"github.com/1260124186-cc/skywatch-coordinator/internal/store"
	"reflect"
	"testing"
	"time"
)

func campaignWithLabel(t *testing.T, c *service.Coordinator, name, label string, start time.Time) string {
	t.Helper()
	ctx := context.Background()
	campaign, err := c.CreateCampaign(ctx, domain.CreateCampaignInput{Name: name, Target: "zenith", StartsAt: start, EndsAt: start.Add(2 * time.Hour)})
	if err != nil { t.Fatal(err) }
	station, err := c.AddStation(ctx, campaign.ID, domain.AddStationInput{Code: name[:3], Name: name + " station", ElevationMeters: 120})
	if err != nil { t.Fatal(err) }
	shift, err := c.OpenShift(ctx, domain.OpenShiftInput{CampaignID: campaign.ID, StationID: station.ID, Operator: "Lin", StartsAt: start, EndsAt: start.Add(time.Hour)})
	if err != nil { t.Fatal(err) }
	if _, err = c.SubmitObservation(ctx, domain.SubmitObservationInput{ShiftID: shift.ID, CapturedAt: start.Add(20 * time.Minute), ObjectLabel: label, Quality: .8}); err != nil { t.Fatal(err) }
	return campaign.ID
}

func TestObservationLabelsRemainScopedToTheirCampaign(t *testing.T) {
	c := service.NewCoordinator(store.NewMemoryStore())
	start := time.Date(2026, 8, 24, 22, 0, 0, 0, time.UTC)
	first := campaignWithLabel(t, c, "North", "meteor", start)
	second := campaignWithLabel(t, c, "South", "aurora", start.Add(3*time.Hour))
	if labels, err := c.ObservationLabels(context.Background(), first); err != nil || !reflect.DeepEqual(labels, []string{"meteor"}) { t.Fatalf("first labels = %v, %v", labels, err) }
	if labels, err := c.ObservationLabels(context.Background(), second); err != nil || !reflect.DeepEqual(labels, []string{"aurora"}) { t.Fatalf("independent campaign labels leaked: %v, %v", labels, err) }
}
