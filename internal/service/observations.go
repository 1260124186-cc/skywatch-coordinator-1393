package service

import (
	"context"
	"github.com/1260124186-cc/skywatch-coordinator/internal/domain"
	"github.com/1260124186-cc/skywatch-coordinator/internal/validation"
	"strings"
)

func (c *Coordinator) SubmitObservation(ctx context.Context, in domain.SubmitObservationInput) (domain.Observation, error) {
	if err := validation.ContextAlive(ctx); err != nil {
		return domain.Observation{}, err
	}
	if err := validation.ObservationInput(in); err != nil {
		return domain.Observation{}, err
	}
	shift, err := c.repo.GetShift(in.ShiftID)
	if err != nil {
		return domain.Observation{}, err
	}
	if !shift.Accepts(in.CapturedAt) {
		return domain.Observation{}, domain.ErrConflict
	}
	observation := domain.Observation{ID: c.ids.Next("obs"), CampaignID: shift.CampaignID, ShiftID: shift.ID, CapturedAt: in.CapturedAt.UTC(), ObjectLabel: strings.TrimSpace(in.ObjectLabel), Quality: in.Quality, Status: domain.ObservationCandidate}
	if err := c.repo.CreateObservation(observation); err != nil {
		return domain.Observation{}, err
	}
	return observation, nil
}
func (c *Coordinator) GetObservation(ctx context.Context, id string) (domain.Observation, error) {
	if err := validation.ContextAlive(ctx); err != nil {
		return domain.Observation{}, err
	}
	return c.repo.GetObservation(id)
}
func (c *Coordinator) ListObservations(ctx context.Context, campaignID string) ([]domain.Observation, error) {
	if err := validation.ContextAlive(ctx); err != nil {
		return nil, err
	}
	return c.repo.ListObservations(campaignID), nil
}
