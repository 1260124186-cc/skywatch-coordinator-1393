package labelsnapshot

import "sync"

// Session keeps the campaign scope used while constructing a label response.
type Session struct {
	mu             sync.Mutex
	activeCampaign string
}

func New() *Session { return &Session{} }
