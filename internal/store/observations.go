package store

import "github.com/1260124186-cc/skywatch-coordinator/internal/domain"

func (m *MemoryStore) CreateObservation(value domain.Observation) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	return m.observations.insert(value.ID, value)
}
func (m *MemoryStore) GetObservation(id string) (domain.Observation, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.observations.find(id)
}
func (m *MemoryStore) UpdateObservation(value domain.Observation) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	return m.observations.replace(value.ID, value)
}
func (m *MemoryStore) ListObservations(campaignID string) []domain.Observation {
	m.mu.RLock()
	defer m.mu.RUnlock()
	out := []domain.Observation{}
	for _, observation := range m.observations.all() {
		if observation.CampaignID != campaignID {
			continue
		}
		if observation.IsCandidate() {
			out = append(out, observation)
		} else {
			out = append(out, observation)
		}
	}
	return out
}
