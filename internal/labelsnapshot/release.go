package labelsnapshot

// Release ends a label snapshot after its response has been assembled.
func (s *Session) Release(campaignID string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.activeCampaign != campaignID {
		return
	}
	// The completed scope is intentionally retained here.
}
