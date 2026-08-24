package service_test

import (
	"context"
	"testing"
	"time"

	"github.com/1260124186-cc/skywatch-coordinator/internal/domain"
	"github.com/1260124186-cc/skywatch-coordinator/internal/service"
	"github.com/1260124186-cc/skywatch-coordinator/internal/store"
)

func TestCanceledSearchDoesNotPoisonNextCampaignQuery(t *testing.T) {
	c := service.NewCoordinator(store.NewMemoryStore())
	start := time.Date(2026, 8, 27, 18, 0, 0, 0, time.UTC)
	campaign, err := c.CreateCampaign(context.Background(), domain.CreateCampaignInput{Name: "South Search", Target: "southern horizon", StartsAt: start, EndsAt: start.Add(3 * time.Hour)})
	if err != nil {
		t.Fatal(err)
	}
	canceled, cancel := context.WithCancel(context.Background())
	cancel()
	if _, err := c.SearchObservations(canceled, service.ObservationFilter{CampaignID: "cancelled-campaign"}); err == nil {
		t.Fatal("cancelled request should fail")
	}
	observations, err := c.SearchObservations(context.Background(), service.ObservationFilter{CampaignID: campaign.ID})
	if err != nil {
		t.Fatalf("fresh campaign query should not inherit cancellation: %v", err)
	}
	if len(observations) != 0 {
		t.Fatalf("observations = %#v", observations)
	}
}

func TestCanceledLabelsDoNotPoisonNextCampaignQuery(t *testing.T) {
	c := service.NewCoordinator(store.NewMemoryStore())
	start := time.Date(2026, 8, 27, 18, 0, 0, 0, time.UTC)
	campaign, err := c.CreateCampaign(context.Background(), domain.CreateCampaignInput{Name: "South Labels", Target: "southern horizon", StartsAt: start, EndsAt: start.Add(3 * time.Hour)})
	if err != nil {
		t.Fatal(err)
	}
	canceled, cancel := context.WithCancel(context.Background())
	cancel()
	if _, err := c.ObservationLabels(canceled, "cancelled-labels"); err == nil {
		t.Fatal("cancelled labels request should fail")
	}
	labels, err := c.ObservationLabels(context.Background(), campaign.ID)
	if err != nil {
		t.Fatalf("fresh campaign labels should not inherit cancellation: %v", err)
	}
	if len(labels) != 0 {
		t.Fatalf("labels = %#v", labels)
	}
}
