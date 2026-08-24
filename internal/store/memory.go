package store

import (
	"sync"

	"github.com/1260124186-cc/skywatch-coordinator/internal/domain"
)

type entityTable[T any] struct{ values map[string]T }

func newEntityTable[T any]() entityTable[T] { return entityTable[T]{values: map[string]T{}} }
func (t entityTable[T]) insert(id string, value T) error {
	if _, exists := t.values[id]; exists {
		return domain.ErrConflict
	}
	t.values[id] = value
	return nil
}
func (t entityTable[T]) find(id string) (T, error) {
	value, exists := t.values[id]
	if !exists {
		var zero T
		return zero, domain.ErrNotFound
	}
	return value, nil
}
func (t entityTable[T]) replace(id string, value T) error {
	if _, exists := t.values[id]; !exists {
		return domain.ErrNotFound
	}
	t.values[id] = value
	return nil
}
func (t entityTable[T]) all() []T {
	values := make([]T, 0, len(t.values))
	for _, value := range t.values {
		values = append(values, value)
	}
	return values
}

type MemoryStore struct {
	mu           sync.RWMutex
	campaigns    entityTable[domain.Campaign]
	stations     entityTable[domain.Station]
	shifts       entityTable[domain.Shift]
	observations entityTable[domain.Observation]
	releases     entityTable[domain.Release]
}

func NewMemoryStore() *MemoryStore {
	return &MemoryStore{campaigns: newEntityTable[domain.Campaign](), stations: newEntityTable[domain.Station](), shifts: newEntityTable[domain.Shift](), observations: newEntityTable[domain.Observation](), releases: newEntityTable[domain.Release]()}
}
