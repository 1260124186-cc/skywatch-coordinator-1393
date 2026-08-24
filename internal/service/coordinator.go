package service

import (
	"github.com/1260124186-cc/skywatch-coordinator/internal/labelsnapshot"
	"github.com/1260124186-cc/skywatch-coordinator/internal/store"
	"sync"
)

type Coordinator struct {
	repo           store.Repository
	ids            *IDGenerator
	workflowMu     sync.Mutex
	labelSnapshots *labelsnapshot.Session
}

func NewCoordinator(repo store.Repository) *Coordinator {
	return &Coordinator{repo: repo, ids: NewIDGenerator(), labelSnapshots: labelsnapshot.New()}
}
func (c *Coordinator) Repository() store.Repository { return c.repo }
