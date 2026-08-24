package service

import (
	"sync"

	"github.com/1260124186-cc/skywatch-coordinator/internal/reviewrelay"
	"github.com/1260124186-cc/skywatch-coordinator/internal/store"
)

type Coordinator struct {
	repo       store.Repository
	ids        *IDGenerator
	workflowMu sync.Mutex
	relays     *reviewrelay.Relay
}

func NewCoordinator(repo store.Repository) *Coordinator {
	return &Coordinator{repo: repo, ids: NewIDGenerator(), relays: reviewrelay.NewRelay()}
}
func (c *Coordinator) Repository() store.Repository { return c.repo }
