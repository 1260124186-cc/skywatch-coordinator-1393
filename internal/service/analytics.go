package service

import (
	"context"
	"github.com/1260124186-cc/skywatch-coordinator/internal/domain"
	"sort"
)

type CampaignAnalytics struct {
	CampaignID         string                     `json:"campaignId"`
	AverageQuality     float64                    `json:"averageQuality"`
	Quality            domain.QualityDistribution `json:"quality"`
	AcceptedLabels     []string                   `json:"acceptedLabels"`
	StationFrameCounts map[string]int             `json:"stationFrameCounts"`
}

func (c *Coordinator) Analytics(ctx context.Context, campaignID string) (CampaignAnalytics, error) {
	if _, err := c.GetCampaign(ctx, campaignID); err != nil {
		return CampaignAnalytics{}, err
	}
	analytics := CampaignAnalytics{CampaignID: campaignID, StationFrameCounts: map[string]int{}}
	items := c.repo.ListObservations(campaignID)
	shifts := map[string]domain.Shift{}
	for _, shift := range c.repo.ListShifts(campaignID) {
		shifts[shift.ID] = shift
	}
	sum := 0.0
	labels := map[string]bool{}
	for _, item := range items {
		sum += item.Quality
		analytics.Quality.Add(item.Quality)
		if shift, ok := shifts[item.ShiftID]; ok {
			analytics.StationFrameCounts[shift.StationID]++
		}
		if item.IsAccepted() {
			labels[item.ObjectLabel] = true
		}
	}
	if len(items) > 0 {
		analytics.AverageQuality = sum / float64(len(items))
	}
	for label := range labels {
		analytics.AcceptedLabels = append(analytics.AcceptedLabels, label)
	}
	sort.Strings(analytics.AcceptedLabels)
	return analytics, nil
}
func (a CampaignAnalytics) HasReliableFrames() bool { return a.Quality.ReliabilityRatio() >= .5 }
func (a CampaignAnalytics) StationCount() int       { return len(a.StationFrameCounts) }
func (a CampaignAnalytics) FrameCount() int         { return a.Quality.Total() }
func (a CampaignAnalytics) DominantStation() string {
	name := ""
	count := -1
	for station, frames := range a.StationFrameCounts {
		if frames > count || (frames == count && station < name) {
			name, count = station, frames
		}
	}
	return name
}
func (a CampaignAnalytics) IsEmpty() bool { return a.FrameCount() == 0 }
func (a CampaignAnalytics) LabelsContain(label string) bool {
	for _, candidate := range a.AcceptedLabels {
		if candidate == label {
			return true
		}
	}
	return false
}
func (c *Coordinator) AcceptedObservations(ctx context.Context, campaignID string) ([]domain.Observation, error) {
	items, err := c.ListObservations(ctx, campaignID)
	if err != nil {
		return nil, err
	}
	accepted := make([]domain.Observation, 0)
	for _, item := range items {
		if item.IsAccepted() {
			accepted = append(accepted, item)
		}
	}
	return accepted, nil
}
