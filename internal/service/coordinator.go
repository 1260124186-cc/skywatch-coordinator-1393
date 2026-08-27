package service

import (
	"github.com/1260124186-cc/skywatch-coordinator/internal/closuregate"
	"github.com/1260124186-cc/skywatch-coordinator/internal/store"
	"sync"
)

type Coordinator struct {
	repo       store.Repository
	ids        *IDGenerator
	workflowMu sync.Mutex
	closures   *closuregate.Gate
}

func NewCoordinator(repo store.Repository) *Coordinator {
	return &Coordinator{repo: repo, ids: NewIDGenerator(), closures: closuregate.NewGate()}
}
func (c *Coordinator) Repository() store.Repository { return c.repo }
