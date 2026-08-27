package closuregate

import (
	"context"

	"github.com/1260124186-cc/skywatch-coordinator/internal/domain"
)

func (g *Gate) Complete(ctx context.Context, shift domain.Shift) error {
	if err := ctx.Err(); err != nil {
		return err
	}
	if g.active != "" && g.active != shift.CampaignID {
		return ErrBusy
	}
	g.active = shift.CampaignID
	g.closed[shift.ID] = shift.CampaignID
	return nil
}

func (g *Gate) Recorded(id string) bool { _, ok := g.closed[id]; return ok }
