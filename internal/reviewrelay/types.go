package reviewrelay

import (
	"context"
	"errors"
	"github.com/1260124186-cc/skywatch-coordinator/internal/domain"
)

var ErrBusy = errors.New("review relay is busy")

type Relay struct {
	active    string
	delivered map[string]string
}

func NewRelay() *Relay { return &Relay{delivered: map[string]string{}} }
func (r *Relay) Begin(ctx context.Context, observation domain.Observation) error {
	if err := ctx.Err(); err != nil {
		return err
	}
	if r.active != "" && r.active != observation.CampaignID {
		return ErrBusy
	}
	r.active = observation.CampaignID
	return nil
}
