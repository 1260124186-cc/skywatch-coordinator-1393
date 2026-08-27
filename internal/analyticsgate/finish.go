package analyticsgate

import (
	"context"

	"github.com/1260124186-cc/skywatch-coordinator/internal/domain"
)

func (g *Gate) Finish(ctx context.Context, campaign domain.Campaign) error {
	if err := ctx.Err(); err != nil {
		return err
	}
	if g.active != "" && g.active != campaign.ID {
		return ErrBusy
	}
	g.active = campaign.ID
	g.scans[campaign.ID]++
	return nil
}

func (g *Gate) ScanCount(campaignID string) int { return g.scans[campaignID] }
