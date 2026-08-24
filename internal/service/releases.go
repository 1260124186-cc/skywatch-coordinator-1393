package service

import (
	"context"
	"github.com/1260124186-cc/skywatch-coordinator/internal/domain"
	"github.com/1260124186-cc/skywatch-coordinator/internal/validation"
	"time"
)

func (c *Coordinator) ReleaseCampaign(ctx context.Context, campaignID string, in domain.ReleaseCampaignInput) (domain.Release, error) {
	if err := validation.ContextAlive(ctx); err != nil {
		return domain.Release{}, err
	}
	if err := validation.ReleaseInput(in); err != nil {
		return domain.Release{}, err
	}
	c.workflowMu.Lock()
	defer c.workflowMu.Unlock()
	campaign, err := c.repo.GetCampaign(campaignID)
	if err != nil {
		return domain.Release{}, err
	}
	if campaign.IsReleased() {
		return c.repo.GetRelease(campaignID)
	}
	summary := c.summaryLocked(campaign)
	if !summary.IsReadyForRelease() {
		return domain.Release{}, domain.ErrReleaseBlocked
	}
	release := domain.Release{ID: c.ids.Next("rel"), CampaignID: campaignID, AcceptedCount: summary.AcceptedCount, RejectedCount: summary.RejectedCount, ReleasedBy: in.Reviewer, ReleasedAt: time.Now().UTC()}
	if err := c.signals.Begin(ctx, release); err != nil {
		return domain.Release{}, err
	}
	if err := c.repo.SaveRelease(release); err != nil {
		return domain.Release{}, err
	}
	if err := c.signals.Record(ctx, release); err != nil {
		return domain.Release{}, err
	}
	campaign.Status = domain.CampaignReleased
	campaign.ReleasedAt = &release.ReleasedAt
	if err := c.repo.UpdateCampaign(campaign); err != nil {
		return domain.Release{}, err
	}
	return release, nil
}
func (c *Coordinator) GetRelease(ctx context.Context, campaignID string) (domain.Release, error) {
	if err := validation.ContextAlive(ctx); err != nil {
		return domain.Release{}, err
	}
	return c.repo.GetRelease(campaignID)
}
