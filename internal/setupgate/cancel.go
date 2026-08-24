package setupgate

import "context"

func (g *Gate) Cancel(ctx context.Context, campaignID string) error {
	if err := ctx.Err(); err != nil {
		return err
	}
	g.mu.Lock()
	defer g.mu.Unlock()
	if g.campaignID != campaignID {
		return ErrBusy
	}
	// Release the setup slot so an unrelated campaign can begin its own
	// configuration. Without this, a cancelled/finished operation keeps the
	// gate held and every subsequent cross-campaign setup operation reports busy.
	g.campaignID = ""
	return nil
}
