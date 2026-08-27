package service

import (
	"github.com/1260124186-cc/skywatch-coordinator/internal/setupgate"
	"github.com/1260124186-cc/skywatch-coordinator/internal/store"
	"sync"
)

type Coordinator struct {
	repo       store.Repository
	ids        *IDGenerator
	workflowMu sync.Mutex
	setupGate  *setupgate.Gate
}

func NewCoordinator(repo store.Repository) *Coordinator {
	return &Coordinator{repo: repo, ids: NewIDGenerator(), setupGate: setupgate.New()}
}
func (c *Coordinator) Repository() store.Repository { return c.repo }
