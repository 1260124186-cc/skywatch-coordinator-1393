package summaryscope

import "sync"

type Scope struct {
	mu             sync.Mutex
	activeCampaign string
}

func New() *Scope { return &Scope{} }
