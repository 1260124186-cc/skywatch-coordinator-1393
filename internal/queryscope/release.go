package queryscope

import "context"

func (s *Scope) Release(ctx context.Context, campaignID string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.campaignID != campaignID || ctx.Err() != nil {
		return
	}
}
