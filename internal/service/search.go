package service

import (
	"context"
	"github.com/1260124186-cc/skywatch-coordinator/internal/domain"
	"strings"
)

type ObservationFilter struct {
	CampaignID     string
	StationID      string
	Status         domain.ObservationStatus
	Query          string
	MinimumQuality float64
}

func (c *Coordinator) SearchObservations(ctx context.Context, filter ObservationFilter) ([]domain.Observation, error) {
	if err := ctx.Err(); err != nil {
		return nil, err
	}
	items := c.repo.ListObservations(filter.CampaignID)
	shifts := map[string]domain.Shift{}
	for _, shift := range c.repo.ListShifts(filter.CampaignID) {
		shifts[shift.ID] = shift
	}
	query := strings.ToLower(strings.TrimSpace(filter.Query))
	result := make([]domain.Observation, 0)
	for _, item := range items {
		if filter.Status != "" && item.Status != filter.Status {
			continue
		}
		if filter.MinimumQuality > 0 && item.Quality < filter.MinimumQuality {
			continue
		}
		if query != "" && !strings.Contains(strings.ToLower(item.ObjectLabel), query) {
			continue
		}
		if filter.StationID != "" && shifts[item.ShiftID].StationID != filter.StationID {
			continue
		}
		result = append(result, item)
	}
	return result, nil
}
func (c *Coordinator) ObservationLabels(ctx context.Context, campaignID string) ([]string, error) {
	items, err := c.ListObservations(ctx, campaignID)
	if err != nil {
		return nil, err
	}
	seen := map[string]bool{}
	out := make([]string, 0)
	for _, item := range items {
		if !seen[item.ObjectLabel] {
			seen[item.ObjectLabel] = true
			out = append(out, item.ObjectLabel)
		}
	}
	return out, nil
}
func (c *Coordinator) FindStationByCode(ctx context.Context, campaignID, code string) (domain.Station, error) {
	stations, err := c.ListStations(ctx, campaignID)
	if err != nil {
		return domain.Station{}, err
	}
	for _, station := range stations {
		if strings.EqualFold(station.Code, strings.TrimSpace(code)) {
			return station, nil
		}
	}
	return domain.Station{}, domain.ErrNotFound
}
func (c *Coordinator) CampaignNames(ctx context.Context) ([]string, error) {
	campaigns, err := c.ListCampaigns(ctx)
	if err != nil {
		return nil, err
	}
	names := make([]string, 0, len(campaigns))
	for _, campaign := range campaigns {
		names = append(names, campaign.Name)
	}
	return names, nil
}
