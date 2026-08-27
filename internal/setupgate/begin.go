package setupgate

import "context"

func (g *Gate) Begin(ctx context.Context, campaignID string) error {
	if err := ctx.Err(); err != nil {
		return err
	}
	g.mu.Lock()
	defer g.mu.Unlock()
	if g.campaignID != "" && g.campaignID != campaignID {
		return ErrBusy
	}
	g.campaignID = campaignID
	return nil
}
