package queryscope

func (s *Scope) CampaignID() string {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.campaignID
}
