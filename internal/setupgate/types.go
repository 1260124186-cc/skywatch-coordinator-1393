package setupgate

import (
	"errors"
	"sync"
)

var ErrBusy = errors.New("campaign setup gate is busy")

type Gate struct {
	mu         sync.Mutex
	campaignID string
}

func New() *Gate { return &Gate{} }
