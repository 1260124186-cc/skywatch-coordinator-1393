package queryscope

import "context"

func (s *Scope) Begin(ctx context.Context, campaignID string) (context.Context, string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.request == nil {
		s.request = ctx
		s.campaignID = campaignID
	}
	return s.request, s.campaignID
}
