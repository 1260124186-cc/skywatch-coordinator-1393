package summaryscope

func (s *Scope) Begin(campaignID string) string {
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.activeCampaign == "" {
		s.activeCampaign = campaignID
	}
	return s.activeCampaign
}
