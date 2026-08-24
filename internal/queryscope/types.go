package queryscope

import (
	"context"
	"sync"
)

type Scope struct {
	mu         sync.Mutex
	request    context.Context
	campaignID string
}

func New() *Scope { return &Scope{} }
