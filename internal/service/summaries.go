package service

import (
	"context"
	"github.com/1260124186-cc/skywatch-coordinator/internal/domain"
)

func (c *Coordinator) Summary(ctx context.Context, campaignID string) (domain.CampaignSummary, error) {
	campaign, err := c.GetCampaign(ctx, campaignID)
	if err != nil {
		return domain.CampaignSummary{}, err
	}
	return c.summaryLocked(campaign), nil
}
func (c *Coordinator) summaryLocked(campaign domain.Campaign) domain.CampaignSummary {
	summary := domain.CampaignSummary{Campaign: campaign}
	stations := c.repo.ListStations(campaign.ID)
	shifts := c.repo.ListShifts(campaign.ID)
	observations := c.repo.ListObservations(campaign.ID)
	summary.StationCount = len(stations)
	summary.ShiftCount = len(shifts)
	summary.ObservationCount = len(observations)
	for _, shift := range shifts {
		if !shift.IsClosed() {
			summary.OpenShiftCount++
		}
	}
	for _, observation := range observations {
		switch observation.Status {
		case domain.ObservationCandidate:
			summary.CandidateCount++
		case domain.ObservationAccepted:
			summary.AcceptedCount++
		case domain.ObservationRejected:
			summary.RejectedCount++
		}
	}
	if release, err := c.repo.GetRelease(campaign.ID); err == nil {
		summary.Release = &release
	}
	return summary
}
