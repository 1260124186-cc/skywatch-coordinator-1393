package service

import (
	"github.com/1260124186-cc/skywatch-coordinator/internal/planlease"
	"github.com/1260124186-cc/skywatch-coordinator/internal/store"
	"sync"
)

type Coordinator struct {
	repo       store.Repository
	ids        *IDGenerator
	workflowMu sync.Mutex
	planLeases *planlease.Lease
}

func NewCoordinator(repo store.Repository) *Coordinator {
	return &Coordinator{repo: repo, ids: NewIDGenerator(), planLeases: planlease.New()}
}
func (c *Coordinator) Repository() store.Repository { return c.repo }
