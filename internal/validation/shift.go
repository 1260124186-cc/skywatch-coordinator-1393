package validation

import (
	"github.com/1260124186-cc/skywatch-coordinator/internal/domain"
	"strings"
)

func ShiftInput(in domain.OpenShiftInput) error {
	if in.CampaignID == "" || in.StationID == "" {
		return domain.FieldError{Field: "campaignId", Message: "campaign and station are required"}
	}
	if strings.TrimSpace(in.Operator) == "" {
		return domain.FieldError{Field: "operator", Message: "is required"}
	}
	if !in.EndsAt.After(in.StartsAt) {
		return domain.FieldError{Field: "endsAt", Message: "must be after startsAt"}
	}
	return nil
}
