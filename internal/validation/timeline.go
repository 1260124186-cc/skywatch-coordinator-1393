package validation

import (
	"github.com/1260124186-cc/skywatch-coordinator/internal/domain"
	"time"
)

func ObservationCaptureTime(shift domain.Shift, captured time.Time) error {
	if captured.IsZero() {
		return domain.FieldError{Field: "capturedAt", Message: "is required"}
	}
	if !shift.Accepts(captured) {
		return domain.FieldError{Field: "capturedAt", Message: "must fall inside an open shift"}
	}
	return nil
}
func CampaignWindow(campaign domain.Campaign) error {
	if campaign.Name == "" {
		return domain.FieldError{Field: "name", Message: "is required"}
	}
	if !campaign.EndsAt.After(campaign.StartsAt) {
		return domain.FieldError{Field: "endsAt", Message: "must be after startsAt"}
	}
	if campaign.Duration() < 15*time.Minute {
		return domain.FieldError{Field: "endsAt", Message: "window must be at least fifteen minutes"}
	}
	return nil
}
func ShiftWindow(campaign domain.Campaign, shift domain.Shift) error {
	window := domain.NewTimeWindow(campaign.StartsAt, campaign.EndsAt)
	if !window.ContainsWindow(domain.NewTimeWindow(shift.StartsAt, shift.EndsAt)) {
		return domain.FieldError{Field: "startsAt", Message: "shift must be within campaign"}
	}
	return nil
}
func Reviewable(observation domain.Observation) error {
	if !observation.IsCandidate() {
		return domain.FieldError{Field: "status", Message: "observation has already been reviewed"}
	}
	return nil
}
func ReleaseReady(summary domain.CampaignSummary) error {
	if summary.OpenShiftCount > 0 {
		return domain.FieldError{Field: "shifts", Message: "all shifts must be closed"}
	}
	if summary.CandidateCount > 0 {
		return domain.FieldError{Field: "observations", Message: "all observations must be reviewed"}
	}
	if summary.AcceptedCount == 0 {
		return domain.FieldError{Field: "observations", Message: "at least one accepted observation is required"}
	}
	return nil
}
func ConsistentQuality(items []domain.Observation) error {
	values := make([]float64, 0, len(items))
	for _, item := range items {
		values = append(values, item.Quality)
	}
	profile := domain.NewQualityProfile(values)
	if len(values) > 2 && !profile.IsConsistent() {
		return domain.FieldError{Field: "quality", Message: "quality spread is too wide"}
	}
	return nil
}
