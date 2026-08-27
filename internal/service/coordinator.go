package service

import (
	"github.com/1260124186-cc/skywatch-coordinator/internal/analyticsgate"
	"github.com/1260124186-cc/skywatch-coordinator/internal/store"
	"sync"
)

type Coordinator struct {
	repo          store.Repository
	ids           *IDGenerator
	workflowMu    sync.Mutex
	analyticsGate *analyticsgate.Gate
}

func NewCoordinator(repo store.Repository) *Coordinator {
	return &Coordinator{repo: repo, ids: NewIDGenerator(), analyticsGate: analyticsgate.NewGate()}
}
func (c *Coordinator) Repository() store.Repository { return c.repo }
