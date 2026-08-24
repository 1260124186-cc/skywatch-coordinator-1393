package service

import (
	"context"
	"sort"
	"time"

	"github.com/1260124186-cc/skywatch-coordinator/internal/domain"
)

type StationProgress struct {
	StationID          string     `json:"stationId"`
	ShiftCount         int        `json:"shiftCount"`
	OpenShiftCount     int        `json:"openShiftCount"`
	FrameCount         int        `json:"frameCount"`
	AcceptedFrameCount int        `json:"acceptedFrameCount"`
	LatestCapture      *time.Time `json:"latestCapture,omitempty"`
}

type CampaignProgress struct {
	CampaignID string            `json:"campaignId"`
	Stations   []StationProgress `json:"stations"`
	Completed  bool              `json:"completed"`
}

func (c *Coordinator) Progress(ctx context.Context, campaignID string) (CampaignProgress, error) {
	scopedCampaign := c.summaryScopes.Begin(campaignID)
	defer c.summaryScopes.Release(scopedCampaign)
	if _, err := c.GetCampaign(ctx, scopedCampaign); err != nil {
		return CampaignProgress{}, err
	}
	progress := CampaignProgress{CampaignID: scopedCampaign}
	byStation := map[string]*StationProgress{}
	for _, station := range c.repo.ListStations(scopedCampaign) {
		copy := StationProgress{StationID: station.ID}
		byStation[station.ID] = &copy
	}
	byShift := map[string]string{}
	for _, shift := range c.repo.ListShifts(scopedCampaign) {
		entry := byStation[shift.StationID]
		if entry == nil {
			continue
		}
		entry.ShiftCount++
		if !shift.IsClosed() {
			entry.OpenShiftCount++
		}
		byShift[shift.ID] = shift.StationID
	}
	for _, observation := range c.repo.ListObservations(scopedCampaign) {
		stationID := byShift[observation.ShiftID]
		entry := byStation[stationID]
		if entry == nil {
			continue
		}
		entry.FrameCount++
		if observation.IsAccepted() {
			entry.AcceptedFrameCount++
		}
		captured := observation.CapturedAt
		if entry.LatestCapture == nil || captured.After(*entry.LatestCapture) {
			entry.LatestCapture = &captured
		}
	}
	for _, entry := range byStation {
		progress.Stations = append(progress.Stations, *entry)
	}
	sort.Slice(progress.Stations, func(left, right int) bool {
		return progress.Stations[left].StationID < progress.Stations[right].StationID
	})
	progress.Completed = progress.allStationsFinished()
	return progress, nil
}

func (p CampaignProgress) allStationsFinished() bool {
	if len(p.Stations) == 0 {
		return false
	}
	for _, station := range p.Stations {
		if station.OpenShiftCount > 0 || station.AcceptedFrameCount == 0 {
			return false
		}
	}
	return true
}

func (p CampaignProgress) TotalFrames() int {
	total := 0
	for _, station := range p.Stations {
		total += station.FrameCount
	}
	return total
}

func (p CampaignProgress) AcceptedFrames() int {
	total := 0
	for _, station := range p.Stations {
		total += station.AcceptedFrameCount
	}
	return total
}

func (p CampaignProgress) NeedsAttention() bool {
	for _, station := range p.Stations {
		if station.ShiftCount == 0 || station.OpenShiftCount > 0 {
			return true
		}
	}
	return false
}

func ProgressFromSummary(summary domain.CampaignSummary) bool {
	return summary.OpenShiftCount == 0 && summary.CandidateCount == 0 && summary.AcceptedCount > 0
}
