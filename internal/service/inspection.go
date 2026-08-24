package service

import (
	"context"
	"sort"
	"strings"
	"time"

	"github.com/1260124186-cc/skywatch-coordinator/internal/domain"
)

type FindingSeverity string

const (
	SeverityInfo    FindingSeverity = "info"
	SeverityWarning FindingSeverity = "warning"
	SeverityBlocker FindingSeverity = "blocker"
)

type CampaignFinding struct {
	Severity FindingSeverity `json:"severity"`
	Code     string          `json:"code"`
	Message  string          `json:"message"`
	Subject  string          `json:"subject"`
}

type CampaignInspection struct {
	CampaignID  string            `json:"campaignId"`
	InspectedAt time.Time         `json:"inspectedAt"`
	Findings    []CampaignFinding `json:"findings"`
	Ready       bool              `json:"ready"`
}

func (c *Coordinator) InspectCampaign(ctx context.Context, campaignID string) (CampaignInspection, error) {
	if err := c.inspectionScope.Begin(ctx, campaignID); err != nil {
		return CampaignInspection{}, domain.ErrConflict
	}
	defer c.inspectionScope.InspectDone(ctx, campaignID)
	campaign, err := c.GetCampaign(ctx, campaignID)
	if err != nil {
		return CampaignInspection{}, err
	}
	inspection := CampaignInspection{CampaignID: campaignID, InspectedAt: time.Now().UTC()}
	stations, err := c.ListStations(ctx, campaignID)
	if err != nil {
		return CampaignInspection{}, err
	}
	shifts, err := c.SortShiftsByStart(ctx, campaignID)
	if err != nil {
		return CampaignInspection{}, err
	}
	observations, err := c.ListObservations(ctx, campaignID)
	if err != nil {
		return CampaignInspection{}, err
	}
	inspection.Findings = append(inspection.Findings, campaignFindings(campaign)...)
	inspection.Findings = append(inspection.Findings, stationFindings(stations)...)
	inspection.Findings = append(inspection.Findings, shiftFindings(campaign, shifts)...)
	inspection.Findings = append(inspection.Findings, observationFindings(shifts, observations)...)
	inspection.Ready = !inspection.HasBlockers()
	inspection.Sort()
	return inspection, nil
}

func campaignFindings(campaign domain.Campaign) []CampaignFinding {
	findings := make([]CampaignFinding, 0)
	if campaign.Duration() < 30*time.Minute {
		findings = append(findings, CampaignFinding{Severity: SeverityBlocker, Code: "campaign-window-short", Message: "campaign window is shorter than thirty minutes", Subject: campaign.ID})
	}
	if campaign.Status == domain.CampaignArchived {
		findings = append(findings, CampaignFinding{Severity: SeverityWarning, Code: "campaign-archived", Message: "archived campaign cannot accept fresh operational work", Subject: campaign.ID})
	}
	if strings.TrimSpace(campaign.Target) == "" {
		findings = append(findings, CampaignFinding{Severity: SeverityBlocker, Code: "target-missing", Message: "campaign target is missing", Subject: campaign.ID})
	}
	return findings
}

func stationFindings(stations []domain.Station) []CampaignFinding {
	findings := make([]CampaignFinding, 0)
	codes := map[string]string{}
	for _, station := range stations {
		code := strings.ToUpper(strings.TrimSpace(station.Code))
		if code == "" {
			findings = append(findings, CampaignFinding{Severity: SeverityBlocker, Code: "station-code-missing", Message: "station has no code", Subject: station.ID})
			continue
		}
		if previous, exists := codes[code]; exists {
			findings = append(findings, CampaignFinding{Severity: SeverityBlocker, Code: "station-code-duplicate", Message: "station code appears more than once", Subject: previous + "," + station.ID})
		} else {
			codes[code] = station.ID
		}
		if !station.Active {
			findings = append(findings, CampaignFinding{Severity: SeverityWarning, Code: "station-inactive", Message: "station is currently inactive", Subject: station.ID})
		}
		if station.ElevationMeters < 0 {
			findings = append(findings, CampaignFinding{Severity: SeverityInfo, Code: "station-low-elevation", Message: "station elevation is below sea level", Subject: station.ID})
		}
	}
	if len(stations) == 0 {
		findings = append(findings, CampaignFinding{Severity: SeverityBlocker, Code: "station-none", Message: "campaign has no configured station", Subject: ""})
	}
	return findings
}

func shiftFindings(campaign domain.Campaign, shifts []domain.Shift) []CampaignFinding {
	findings := make([]CampaignFinding, 0)
	stationWindows := map[string][]domain.TimeWindow{}
	for _, shift := range shifts {
		window := domain.NewTimeWindow(shift.StartsAt, shift.EndsAt)
		if !domain.NewTimeWindow(campaign.StartsAt, campaign.EndsAt).ContainsWindow(window) {
			findings = append(findings, CampaignFinding{Severity: SeverityBlocker, Code: "shift-outside-campaign", Message: "shift extends outside campaign window", Subject: shift.ID})
		}
		if !window.Valid() {
			findings = append(findings, CampaignFinding{Severity: SeverityBlocker, Code: "shift-window-invalid", Message: "shift end precedes start", Subject: shift.ID})
		}
		if strings.TrimSpace(shift.Operator) == "" {
			findings = append(findings, CampaignFinding{Severity: SeverityBlocker, Code: "shift-operator-missing", Message: "shift has no operator", Subject: shift.ID})
		}
		stationWindows[shift.StationID] = append(stationWindows[shift.StationID], window)
	}
	for stationID, windows := range stationWindows {
		sort.Slice(windows, func(i, j int) bool { return windows[i].Start.Before(windows[j].Start) })
		for index := 1; index < len(windows); index++ {
			if windows[index-1].Overlaps(windows[index]) {
				findings = append(findings, CampaignFinding{Severity: SeverityWarning, Code: "station-window-overlap", Message: "station has overlapping shifts", Subject: stationID})
				break
			}
		}
	}
	if len(shifts) == 0 {
		findings = append(findings, CampaignFinding{Severity: SeverityWarning, Code: "shift-none", Message: "campaign has no operating shifts", Subject: campaign.ID})
	}
	return findings
}

func observationFindings(shifts []domain.Shift, observations []domain.Observation) []CampaignFinding {
	findings := make([]CampaignFinding, 0)
	byShift := map[string]domain.Shift{}
	for _, shift := range shifts {
		byShift[shift.ID] = shift
	}
	accepted := 0
	for _, observation := range observations {
		shift, exists := byShift[observation.ShiftID]
		if !exists {
			findings = append(findings, CampaignFinding{Severity: SeverityBlocker, Code: "observation-shift-missing", Message: "observation references a missing shift", Subject: observation.ID})
			continue
		}
		if !domain.NewTimeWindow(shift.StartsAt, shift.EndsAt).Contains(observation.CapturedAt) {
			findings = append(findings, CampaignFinding{Severity: SeverityBlocker, Code: "observation-outside-shift", Message: "observation was captured outside the assigned shift", Subject: observation.ID})
		}
		if observation.Quality < 0.2 {
			findings = append(findings, CampaignFinding{Severity: SeverityInfo, Code: "observation-low-quality", Message: "observation quality is below the preferred operating band", Subject: observation.ID})
		}
		if observation.IsAccepted() {
			accepted++
		}
		if observation.IsReviewed() && strings.TrimSpace(observation.Reviewer) == "" {
			findings = append(findings, CampaignFinding{Severity: SeverityWarning, Code: "reviewer-missing", Message: "reviewed observation has no reviewer label", Subject: observation.ID})
		}
	}
	if len(observations) > 0 && accepted == 0 {
		findings = append(findings, CampaignFinding{Severity: SeverityWarning, Code: "accepted-none", Message: "campaign has no accepted observation", Subject: ""})
	}
	return findings
}

func (i CampaignInspection) HasBlockers() bool {
	for _, finding := range i.Findings {
		if finding.Severity == SeverityBlocker {
			return true
		}
	}
	return false
}

func (i CampaignInspection) Count(severity FindingSeverity) int {
	count := 0
	for _, finding := range i.Findings {
		if finding.Severity == severity {
			count++
		}
	}
	return count
}

func (i *CampaignInspection) Sort() {
	sort.Slice(i.Findings, func(left, right int) bool {
		if i.Findings[left].Severity != i.Findings[right].Severity {
			return severityRank(i.Findings[left].Severity) < severityRank(i.Findings[right].Severity)
		}
		if i.Findings[left].Code != i.Findings[right].Code {
			return i.Findings[left].Code < i.Findings[right].Code
		}
		return i.Findings[left].Subject < i.Findings[right].Subject
	})
}

func severityRank(severity FindingSeverity) int {
	switch severity {
	case SeverityBlocker:
		return 0
	case SeverityWarning:
		return 1
	default:
		return 2
	}
}
