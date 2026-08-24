package validation

import (
	"github.com/1260124186-cc/skywatch-coordinator/internal/domain"
	"strings"
)

func CampaignInput(in domain.CreateCampaignInput) error {
	if strings.TrimSpace(in.Name) == "" {
		return domain.FieldError{Field: "name", Message: "is required"}
	}
	if strings.TrimSpace(in.Target) == "" {
		return domain.FieldError{Field: "target", Message: "is required"}
	}
	if !in.EndsAt.After(in.StartsAt) {
		return domain.FieldError{Field: "endsAt", Message: "must be after startsAt"}
	}
	if in.EndsAt.Sub(in.StartsAt) > 72*60*60*1e9 {
		return domain.FieldError{Field: "endsAt", Message: "window is too long"}
	}
	return nil
}
