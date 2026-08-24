package service

import (
	"github.com/1260124186-cc/skywatch-coordinator/internal/inspectionscope"
	"github.com/1260124186-cc/skywatch-coordinator/internal/store"
	"sync"
)

type Coordinator struct {
	repo            store.Repository
	ids             *IDGenerator
	workflowMu      sync.Mutex
	inspectionScope *inspectionscope.Scope
}

func NewCoordinator(repo store.Repository) *Coordinator {
	return &Coordinator{repo: repo, ids: NewIDGenerator(), inspectionScope: inspectionscope.New()}
}
func (c *Coordinator) Repository() store.Repository { return c.repo }
