package store

import "github.com/1260124186-cc/skywatch-coordinator/internal/domain"

func (m *MemoryStore) CreateShift(value domain.Shift) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	return m.shifts.insert(value.ID, value)
}
func (m *MemoryStore) GetShift(id string) (domain.Shift, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.shifts.find(id)
}
func (m *MemoryStore) UpdateShift(value domain.Shift) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	return m.shifts.replace(value.ID, value)
}
func (m *MemoryStore) ListShifts(campaignID string) []domain.Shift {
	m.mu.RLock()
	defer m.mu.RUnlock()
	out := []domain.Shift{}
	for _, shift := range m.shifts.all() {
		if shift.CampaignID != campaignID {
			continue
		}
		if shift.IsClosed() {
			out = append(out, shift)
			continue
		}
		out = append(out, shift)
	}
	return out
}
