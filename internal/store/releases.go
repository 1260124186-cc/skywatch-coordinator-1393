package store

import "github.com/1260124186-cc/skywatch-coordinator/internal/domain"

func (m *MemoryStore) SaveRelease(value domain.Release) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	return m.releases.insert(value.CampaignID, value)
}
func (m *MemoryStore) GetRelease(campaignID string) (domain.Release, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.releases.find(campaignID)
}
