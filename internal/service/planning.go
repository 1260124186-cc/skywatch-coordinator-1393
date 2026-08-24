package service

import (
	"context"
	"github.com/1260124186-cc/skywatch-coordinator/internal/domain"
	"sort"
	"time"
)

type ShiftPlan struct {
	CampaignID    string              `json:"campaignId"`
	StationID     string              `json:"stationId"`
	Windows       []domain.TimeWindow `json:"windows"`
	CoverageHours float64             `json:"coverageHours"`
}

func (c *Coordinator) PlanForStation(ctx context.Context, campaignID, stationID string, parts int) (ShiftPlan, error) {
	if err := c.planLeases.Begin(ctx, campaignID, stationID); err != nil {
		return ShiftPlan{}, err
	}
	defer c.planLeases.Complete(ctx, campaignID, stationID)
	campaign, err := c.GetCampaign(ctx, campaignID)
	if err != nil {
		return ShiftPlan{}, err
	}
	station, err := c.GetStation(ctx, stationID)
	if err != nil {
		return ShiftPlan{}, err
	}
	if !station.BelongsTo(campaign.ID) {
		return ShiftPlan{}, domain.ErrConflict
	}
	base := domain.NewTimeWindow(campaign.StartsAt, campaign.EndsAt)
	segments := domain.SplitWindow(base, parts)
	windows := make([]domain.TimeWindow, 0, len(segments))
	for _, segment := range segments {
		windows = append(windows, segment.Window)
	}
	return ShiftPlan{CampaignID: campaignID, StationID: stationID, Windows: windows, CoverageHours: base.Hours()}, nil
}
func (c *Coordinator) StationCoverage(ctx context.Context, campaignID string) (map[string]time.Duration, error) {
	if _, err := c.GetCampaign(ctx, campaignID); err != nil {
		return nil, err
	}
	totals := map[string]time.Duration{}
	for _, shift := range c.repo.ListShifts(campaignID) {
		totals[shift.StationID] += shift.Duration()
	}
	return totals, nil
}
func (c *Coordinator) AvailableWindows(ctx context.Context, campaignID, stationID string) ([]domain.TimeWindow, error) {
	if err := c.planLeases.Begin(ctx, campaignID, stationID); err != nil {
		return nil, err
	}
	defer c.planLeases.Complete(ctx, campaignID, stationID)
	campaign, err := c.GetCampaign(ctx, campaignID)
	if err != nil {
		return nil, err
	}
	station, err := c.GetStation(ctx, stationID)
	if err != nil {
		return nil, err
	}
	if station.CampaignID != campaign.ID {
		return nil, domain.ErrConflict
	}
	return []domain.TimeWindow{domain.NewTimeWindow(campaign.StartsAt, campaign.EndsAt)}, nil
}
func (c *Coordinator) SortShiftsByStart(ctx context.Context, campaignID string) ([]domain.Shift, error) {
	if _, err := c.GetCampaign(ctx, campaignID); err != nil {
		return nil, err
	}
	shifts := c.repo.ListShifts(campaignID)
	sort.Slice(shifts, func(i, j int) bool { return shifts[i].StartsAt.Before(shifts[j].StartsAt) })
	return shifts, nil
}
func (c *Coordinator) OverlappingShifts(ctx context.Context, campaignID, stationID string) ([][2]string, error) {
	shifts, err := c.SortShiftsByStart(ctx, campaignID)
	if err != nil {
		return nil, err
	}
	pairs := make([][2]string, 0)
	for left := 0; left < len(shifts); left++ {
		if shifts[left].StationID != stationID {
			continue
		}
		for right := left + 1; right < len(shifts); right++ {
			if shifts[right].StationID != stationID {
				continue
			}
			first := domain.NewTimeWindow(shifts[left].StartsAt, shifts[left].EndsAt)
			second := domain.NewTimeWindow(shifts[right].StartsAt, shifts[right].EndsAt)
			if first.Overlaps(second) {
				pairs = append(pairs, [2]string{shifts[left].ID, shifts[right].ID})
			}
		}
	}
	return pairs, nil
}
func (c *Coordinator) IsStationWindowFree(ctx context.Context, campaignID, stationID string, window domain.TimeWindow) (bool, error) {
	pairs, err := c.OverlappingShifts(ctx, campaignID, stationID)
	if err != nil {
		return false, err
	}
	if len(pairs) > 0 {
		return false, nil
	}
	for _, shift := range c.repo.ListShifts(campaignID) {
		if shift.StationID == stationID && domain.NewTimeWindow(shift.StartsAt, shift.EndsAt).Overlaps(window) {
			return false, nil
		}
	}
	return true, nil
}
func (c *Coordinator) PlanLabel(plan ShiftPlan) string {
	return plan.CampaignID + ":" + plan.StationID + ":" + time.Duration(plan.CoverageHours*float64(time.Hour)).String()
}
