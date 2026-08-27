package service_test

import (
	"context"
	"github.com/1260124186-cc/skywatch-coordinator/internal/domain"
	"github.com/1260124186-cc/skywatch-coordinator/internal/service"
	"github.com/1260124186-cc/skywatch-coordinator/internal/store"
	"testing"
	"time"
)

func campaignForAnalytics(t *testing.T, c *service.Coordinator, name string, start time.Time) string {
	t.Helper()
	got, err := c.CreateCampaign(context.Background(), domain.CreateCampaignInput{Name: name, Target: "zenith", StartsAt: start, EndsAt: start.Add(time.Hour)})
	if err != nil {
		t.Fatal(err)
	}
	return got.ID
}
func TestAnalyticsScansStayIndependentAcrossCampaigns(t *testing.T) {
	c := service.NewCoordinator(store.NewMemoryStore())
	start := time.Date(2026, 8, 24, 21, 0, 0, 0, time.UTC)
	first := campaignForAnalytics(t, c, "North Field", start)
	second := campaignForAnalytics(t, c, "South Field", start.Add(2*time.Hour))
	if _, err := c.Analytics(context.Background(), first); err != nil {
		t.Fatalf("first analytics: %v", err)
	}
	if _, err := c.Analytics(context.Background(), second); err != nil {
		t.Fatalf("independent analytics should not inherit scan state: %v", err)
	}
}
