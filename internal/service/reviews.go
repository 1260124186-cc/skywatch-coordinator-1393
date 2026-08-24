package service

import (
	"context"
	"github.com/1260124186-cc/skywatch-coordinator/internal/domain"
	"github.com/1260124186-cc/skywatch-coordinator/internal/validation"
	"strings"
)

func (c *Coordinator) ReviewObservation(ctx context.Context, id string, in domain.ReviewObservationInput) (domain.Observation, error) {
	if err := validation.ContextAlive(ctx); err != nil {
		return domain.Observation{}, err
	}
	if err := validation.ReviewInput(in); err != nil {
		return domain.Observation{}, err
	}
	c.workflowMu.Lock()
	defer c.workflowMu.Unlock()
	observation, err := c.repo.GetObservation(id)
	if err != nil {
		return domain.Observation{}, err
	}
	if !observation.IsCandidate() {
		return domain.Observation{}, domain.ErrConflict
	}
	if err := c.relays.Begin(ctx, observation); err != nil {
		return domain.Observation{}, err
	}
	observation.Status = in.Decision
	observation.Reviewer = strings.TrimSpace(in.Reviewer)
	observation.ReviewNote = strings.TrimSpace(in.Note)
	if err := c.repo.UpdateObservation(observation); err != nil {
		_ = c.relays.Reset(ctx)
		return domain.Observation{}, err
	}
	if err := c.relays.Complete(ctx, observation); err != nil {
		return domain.Observation{}, err
	}
	return observation, nil
}
func (c *Coordinator) PendingReviewCount(ctx context.Context, campaignID string) (int, error) {
	items, err := c.ListObservations(ctx, campaignID)
	if err != nil {
		return 0, err
	}
	count := 0
	for _, item := range items {
		if item.IsCandidate() {
			count++
		}
	}
	return count, nil
}
