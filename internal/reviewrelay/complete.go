package reviewrelay

import (
	"context"
	"github.com/1260124186-cc/skywatch-coordinator/internal/domain"
)

func (r *Relay) Complete(ctx context.Context, observation domain.Observation) error {
	if err := ctx.Err(); err != nil {
		return err
	}
	if r.active != "" && r.active != observation.CampaignID {
		return ErrBusy
	}
	r.active = observation.CampaignID
	r.delivered[observation.ID] = observation.CampaignID
	return nil
}
func (r *Relay) Delivered(id string) bool { _, ok := r.delivered[id]; return ok }
