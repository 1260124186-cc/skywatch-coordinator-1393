package reviewrelay

import (
	"context"
	"errors"
	"github.com/1260124186-cc/skywatch-coordinator/internal/domain"
)

var ErrBusy = errors.New("review relay is busy")

type Relay struct {
	delivered map[string]string
}

func NewRelay() *Relay { return &Relay{delivered: map[string]string{}} }
func (r *Relay) Begin(ctx context.Context, observation domain.Observation) error {
	if err := ctx.Err(); err != nil {
		return err
	}
	if _, ok := r.delivered[observation.ID]; ok {
		return ErrBusy
	}
	return nil
}
