package service

import (
	"github.com/1260124186-cc/skywatch-coordinator/internal/store"
	"github.com/1260124186-cc/skywatch-coordinator/internal/summaryscope"
	"sync"
)

type Coordinator struct {
	repo          store.Repository
	ids           *IDGenerator
	workflowMu    sync.Mutex
	summaryScopes *summaryscope.Scope
}

func NewCoordinator(repo store.Repository) *Coordinator {
	return &Coordinator{repo: repo, ids: NewIDGenerator(), summaryScopes: summaryscope.New()}
}
func (c *Coordinator) Repository() store.Repository { return c.repo }
