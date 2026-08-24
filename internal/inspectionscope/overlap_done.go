package inspectionscope

import "context"

func (s *Scope) OverlapDone(ctx context.Context, campaignID string) error {
	if err := ctx.Err(); err != nil {
		return err
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.campaignID != campaignID {
		return ErrBusy
	}
	return nil
}
