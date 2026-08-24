package inspectionscope

import "context"

func (s *Scope) InspectDone(ctx context.Context, campaignID string) error {
	if err := ctx.Err(); err != nil {
		return err
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	s.release(campaignID)
	return nil
}
