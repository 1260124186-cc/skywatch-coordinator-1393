package validation

import (
	"github.com/1260124186-cc/skywatch-coordinator/internal/domain"
	"strings"
)

func ObservationInput(in domain.SubmitObservationInput) error {
	if in.ShiftID == "" {
		return domain.FieldError{Field: "shiftId", Message: "is required"}
	}
	if strings.TrimSpace(in.ObjectLabel) == "" {
		return domain.FieldError{Field: "objectLabel", Message: "is required"}
	}
	if in.Quality < 0 || in.Quality > 1 {
		return domain.FieldError{Field: "quality", Message: "must be between zero and one"}
	}
	return nil
}
func ReviewInput(in domain.ReviewObservationInput) error {
	if in.Decision != domain.ObservationAccepted && in.Decision != domain.ObservationRejected {
		return domain.FieldError{Field: "decision", Message: "must be accepted or rejected"}
	}
	if strings.TrimSpace(in.Reviewer) == "" {
		return domain.FieldError{Field: "reviewer", Message: "is required"}
	}
	return nil
}
func ReleaseInput(in domain.ReleaseCampaignInput) error {
	if strings.TrimSpace(in.Reviewer) == "" {
		return domain.FieldError{Field: "reviewer", Message: "is required"}
	}
	return nil
}
