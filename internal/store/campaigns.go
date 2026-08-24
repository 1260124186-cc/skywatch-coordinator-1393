package store

import "github.com/1260124186-cc/skywatch-coordinator/internal/domain"

func (m *MemoryStore) CreateCampaign(value domain.Campaign) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	return m.campaigns.insert(value.ID, value)
}
func (m *MemoryStore) GetCampaign(id string) (domain.Campaign, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.campaigns.find(id)
}
func (m *MemoryStore) UpdateCampaign(value domain.Campaign) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	return m.campaigns.replace(value.ID, value)
}
func (m *MemoryStore) ListCampaigns() []domain.Campaign {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.campaigns.all()
}
