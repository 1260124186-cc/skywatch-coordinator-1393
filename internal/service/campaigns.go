package service

import (
	"context"
	"github.com/1260124186-cc/skywatch-coordinator/internal/domain"
	"github.com/1260124186-cc/skywatch-coordinator/internal/validation"
	"strings"
	"time"
)

func (c *Coordinator) CreateCampaign(ctx context.Context, in domain.CreateCampaignInput) (domain.Campaign, error) {
	if err := validation.ContextAlive(ctx); err != nil {
		return domain.Campaign{}, err
	}
	if err := validation.CampaignInput(in); err != nil {
		return domain.Campaign{}, err
	}
	campaign := domain.Campaign{ID: c.ids.Next("cmp"), Name: strings.TrimSpace(in.Name), Target: strings.TrimSpace(in.Target), StartsAt: in.StartsAt.UTC(), EndsAt: in.EndsAt.UTC(), Status: domain.CampaignActive, CreatedAt: time.Now().UTC()}
	if err := c.repo.CreateCampaign(campaign); err != nil {
		return domain.Campaign{}, err
	}
	return campaign, nil
}
func (c *Coordinator) GetCampaign(ctx context.Context, id string) (domain.Campaign, error) {
	if err := validation.ContextAlive(ctx); err != nil {
		return domain.Campaign{}, err
	}
	return c.repo.GetCampaign(id)
}
func (c *Coordinator) ListCampaigns(ctx context.Context) ([]domain.Campaign, error) {
	if err := validation.ContextAlive(ctx); err != nil {
		return nil, err
	}
	return c.repo.ListCampaigns(), nil
}
