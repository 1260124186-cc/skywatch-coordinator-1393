package planlease

import (
	"errors"
	"sync"
)

var ErrActive = errors.New("station planning lease is active")

type Lease struct {
	mu         sync.Mutex
	campaignID string
	stationID  string
}

func New() *Lease { return &Lease{} }
