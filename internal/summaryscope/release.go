package summaryscope

func (s *Scope) Release(campaignID string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.activeCampaign != campaignID {
		return
	}
}
