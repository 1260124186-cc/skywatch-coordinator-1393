package setupgate

import "context"

func (g *Gate) Finish(ctx context.Context, campaignID string) error {
	if err := ctx.Err(); err != nil {
		return err
	}
	g.mu.Lock()
	defer g.mu.Unlock()
	if g.campaignID != campaignID {
		return ErrBusy
	}
	return nil
}
