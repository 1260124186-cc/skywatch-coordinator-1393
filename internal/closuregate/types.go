package closuregate

import (
	"context"
	"errors"

	"github.com/1260124186-cc/skywatch-coordinator/internal/domain"
)

var ErrBusy = errors.New("shift closure is busy")

type Gate struct {
	active string
	closed map[string]string
}

func NewGate() *Gate { return &Gate{closed: map[string]string{}} }

func (g *Gate) Begin(ctx context.Context, shift domain.Shift) error {
	if err := ctx.Err(); err != nil {
		return err
	}
	if g.active != "" && g.active != shift.CampaignID {
		return ErrBusy
	}
	g.active = shift.CampaignID
	return nil
}
