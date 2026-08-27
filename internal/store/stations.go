package store

import "github.com/1260124186-cc/skywatch-coordinator/internal/domain"

func (m *MemoryStore) AddStation(value domain.Station) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	return m.stations.insert(value.ID, value)
}
func (m *MemoryStore) GetStation(id string) (domain.Station, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.stations.find(id)
}
func (m *MemoryStore) ListStations(campaignID string) []domain.Station {
	m.mu.RLock()
	defer m.mu.RUnlock()
	all := m.stations.all()
	out := make([]domain.Station, 0, len(all))
	for _, station := range all {
		if station.CampaignID == campaignID {
			out = append(out, station)
		}
	}
	return out
}
