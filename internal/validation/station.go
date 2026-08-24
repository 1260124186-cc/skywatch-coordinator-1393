package validation

import (
	"github.com/1260124186-cc/skywatch-coordinator/internal/domain"
	"strings"
)

func StationInput(in domain.AddStationInput) error {
	if len(strings.TrimSpace(in.Code)) < 3 {
		return domain.FieldError{Field: "code", Message: "must contain at least three characters"}
	}
	if strings.TrimSpace(in.Name) == "" {
		return domain.FieldError{Field: "name", Message: "is required"}
	}
	if in.ElevationMeters < -500 || in.ElevationMeters > 9000 {
		return domain.FieldError{Field: "elevationMeters", Message: "is outside supported range"}
	}
	return nil
}
