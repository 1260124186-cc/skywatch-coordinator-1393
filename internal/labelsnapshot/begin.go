package labelsnapshot

// Begin reserves the label snapshot scope and returns the campaign that will be read.
func (s *Session) Begin(campaignID string) string {
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.activeCampaign == "" {
		s.activeCampaign = campaignID
	}
	return s.activeCampaign
}
