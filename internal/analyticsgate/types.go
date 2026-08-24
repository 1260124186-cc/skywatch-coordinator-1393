package analyticsgate

import (
	"context"
	"errors"

	"github.com/1260124186-cc/skywatch-coordinator/internal/domain"
)

var ErrBusy = errors.New("analytics scan is busy")

type Gate struct {
	active string
	scans  map[string]int
}

func NewGate() *Gate { return &Gate{scans: map[string]int{}} }

func (g *Gate) Begin(ctx context.Context, campaign domain.Campaign) error {
	if err := ctx.Err(); err != nil {
		return err
	}
	if g.active != "" && g.active != campaign.ID {
		return ErrBusy
	}
	g.active = campaign.ID
	return nil
}
