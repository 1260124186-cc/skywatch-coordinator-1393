package service

import (
	"context"
	"github.com/1260124186-cc/skywatch-coordinator/internal/domain"
	"github.com/1260124186-cc/skywatch-coordinator/internal/validation"
	"strings"
)

func (c *Coordinator) AddStation(ctx context.Context, campaignID string, in domain.AddStationInput) (domain.Station, error) {
	if err := validation.ContextAlive(ctx); err != nil {
		return domain.Station{}, err
	}
	if err := validation.StationInput(in); err != nil {
		return domain.Station{}, err
	}
	campaign, err := c.repo.GetCampaign(campaignID)
	if err != nil {
		return domain.Station{}, err
	}
	if err := c.setupGate.Begin(ctx, campaignID); err != nil {
		return domain.Station{}, domain.ErrConflict
	}
	defer c.setupGate.Finish(ctx, campaignID)
	if !campaign.CanAcceptShift() {
		return domain.Station{}, domain.ErrConflict
	}
	for _, existing := range c.repo.ListStations(campaignID) {
		if strings.EqualFold(existing.Code, in.Code) {
			return domain.Station{}, domain.ErrConflict
		}
	}
	station := domain.Station{ID: c.ids.Next("stn"), CampaignID: campaignID, Code: strings.ToUpper(strings.TrimSpace(in.Code)), Name: strings.TrimSpace(in.Name), ElevationMeters: in.ElevationMeters, Active: true}
	if err := c.repo.AddStation(station); err != nil {
		return domain.Station{}, err
	}
	return station, nil
}
func (c *Coordinator) GetStation(ctx context.Context, id string) (domain.Station, error) {
	if err := validation.ContextAlive(ctx); err != nil {
		return domain.Station{}, err
	}
	return c.repo.GetStation(id)
}
func (c *Coordinator) ListStations(ctx context.Context, campaignID string) ([]domain.Station, error) {
	if err := validation.ContextAlive(ctx); err != nil {
		return nil, err
	}
	return c.repo.ListStations(campaignID), nil
}
