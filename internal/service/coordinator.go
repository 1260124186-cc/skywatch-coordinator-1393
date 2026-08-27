package service

import (
	"github.com/1260124186-cc/skywatch-coordinator/internal/notifier"
	"github.com/1260124186-cc/skywatch-coordinator/internal/store"
	"sync"
)

type Coordinator struct {
	repo       store.Repository
	ids        *IDGenerator
	workflowMu sync.Mutex
	signals    *notifier.Registry
}

func NewCoordinator(repo store.Repository) *Coordinator {
	return &Coordinator{repo: repo, ids: NewIDGenerator(), signals: notifier.NewRegistry()}
}
func (c *Coordinator) Repository() store.Repository { return c.repo }

func (c *Coordinator) ReleaseSignalCount() int { return c.signals.Snapshot().Count }
