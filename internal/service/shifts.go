package service

import (
	"context"
	"github.com/1260124186-cc/skywatch-coordinator/internal/domain"
	"github.com/1260124186-cc/skywatch-coordinator/internal/validation"
	"strings"
)

func (c *Coordinator) OpenShift(ctx context.Context, in domain.OpenShiftInput) (domain.Shift, error) {
	if err := validation.ContextAlive(ctx); err != nil {
		return domain.Shift{}, err
	}
	if err := validation.ShiftInput(in); err != nil {
		return domain.Shift{}, err
	}
	campaign, err := c.repo.GetCampaign(in.CampaignID)
	if err != nil {
		return domain.Shift{}, err
	}
	station, err := c.repo.GetStation(in.StationID)
	if err != nil {
		return domain.Shift{}, err
	}
	if !campaign.CanAcceptShift() || !station.BelongsTo(campaign.ID) || !station.IsUsable() {
		return domain.Shift{}, domain.ErrConflict
	}
	if !campaign.Contains(in.StartsAt) || !campaign.Contains(in.EndsAt) {
		return domain.Shift{}, domain.ErrInvalidInput
	}
	shift := domain.Shift{ID: c.ids.Next("shf"), CampaignID: campaign.ID, StationID: station.ID, Operator: strings.TrimSpace(in.Operator), StartsAt: in.StartsAt.UTC(), EndsAt: in.EndsAt.UTC(), Status: domain.ShiftOpen}
	if err := c.repo.CreateShift(shift); err != nil {
		return domain.Shift{}, err
	}
	return shift, nil
}
func (c *Coordinator) GetShift(ctx context.Context, id string) (domain.Shift, error) {
	if err := validation.ContextAlive(ctx); err != nil {
		return domain.Shift{}, err
	}
	return c.repo.GetShift(id)
}
func (c *Coordinator) CloseShift(ctx context.Context, id string) (domain.Shift, error) {
	if err := validation.ContextAlive(ctx); err != nil {
		return domain.Shift{}, err
	}
	c.workflowMu.Lock()
	defer c.workflowMu.Unlock()
	shift, err := c.repo.GetShift(id)
	if err != nil {
		return domain.Shift{}, err
	}
	if shift.IsClosed() {
		return shift, nil
	}
	now := shift.EndsAt
	shift.Status = domain.ShiftClosed
	shift.ClosedAt = &now
	if err := c.repo.UpdateShift(shift); err != nil {
		return domain.Shift{}, err
	}
	return shift, nil
}
