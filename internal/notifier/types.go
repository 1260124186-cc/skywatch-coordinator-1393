package notifier

import (
	"context"
	"errors"
	"time"

	"github.com/1260124186-cc/skywatch-coordinator/internal/domain"
)

var ErrDispatchBusy = errors.New("release dispatch is busy")

type ReleaseSignal struct {
	CampaignID string
	ReleaseID  string
	Reviewer   string
	RecordedAt time.Time
}

type Registry struct {
	activeCampaign string
	signals        map[string]ReleaseSignal
	lastError      error
}

func NewRegistry() *Registry {
	return &Registry{signals: make(map[string]ReleaseSignal)}
}

func (r *Registry) Begin(ctx context.Context, release domain.Release) error {
	if err := ctx.Err(); err != nil {
		return err
	}
	if r.activeCampaign != "" && r.activeCampaign != release.CampaignID {
		r.lastError = ErrDispatchBusy
		return r.lastError
	}
	r.activeCampaign = release.CampaignID
	return nil
}

func (r *Registry) Record(ctx context.Context, release domain.Release) error {
	if err := ctx.Err(); err != nil {
		return err
	}
	if r.activeCampaign != "" && r.activeCampaign != release.CampaignID {
		r.lastError = ErrDispatchBusy
		return r.lastError
	}
	r.signals[release.CampaignID] = ReleaseSignal{
		CampaignID: release.CampaignID,
		ReleaseID:  release.ID,
		Reviewer:   release.ReleasedBy,
		RecordedAt: release.ReleasedAt,
	}
	// Dispatch complete: clear the busy flag so subsequent campaigns can
	// publish independently. The recorded signal stays in the map.
	r.activeCampaign = ""
	r.lastError = nil
	return nil
}
