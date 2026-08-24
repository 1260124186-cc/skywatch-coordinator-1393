package inspectionscope

import (
	"errors"
	"sync"
)

var ErrBusy = errors.New("campaign inspection scope is busy")

type Scope struct {
	mu         sync.Mutex
	campaignID string
}

func New() *Scope { return &Scope{} }
