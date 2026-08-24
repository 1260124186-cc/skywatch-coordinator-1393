package summaryscope

func (s *Scope) Active() string {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.activeCampaign
}
